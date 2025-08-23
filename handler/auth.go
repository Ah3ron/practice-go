package handler

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"app/config"
	"app/database"
	"app/model"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error on login request",
			"errors":  err.Error(),
		})
	}

	username := strings.TrimSpace(input.Username)
	pass := input.Password

	var (
		userModel *model.User
		err       error
	)

	if valid(username) {
		userModel, err = getUserByEmail(username)
	} else {
		userModel, err = getUserByUsername(username)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Internal Server Error",
		})
	}
	if userModel == nil || !CheckPasswordHash(pass, userModel.Password) {
		// не раскрываем детали
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid username or password",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userModel.Username
	claims["user_id"] = userModel.ID
	claims["exp"] = time.Now().Add(72 * time.Hour).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Return both user and token as expected by frontend
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Success login",
		"data": fiber.Map{
			"user": fiber.Map{
				"id":         userModel.ID,
				"username":   userModel.Username,
				"email":      userModel.Email,
				"names":      userModel.Names,
				"created_at": userModel.CreatedAt,
				"updated_at": userModel.UpdatedAt,
			},
			"token": t,
		},
	})
}

func Register(c *fiber.Ctx) error {
	var u model.User
	if err := c.BodyParser(&u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	hash, _ := hashPassword(u.Password)
	u.Password = hash
	if err := database.DB.Create(&u).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "data": map[string]interface{}{
		"id": u.ID, "username": u.Username, "email": u.Email,
	}})
}
