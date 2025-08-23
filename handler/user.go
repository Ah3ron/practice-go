package handler

import (
	"strconv"
	"time"

	"app/database"
	"app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	db := database.DB
	var user model.User
	db.First(&user, id)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(p, user.Password) {
		return false
	}
	return true
}

// GetAllUsers get all users
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []model.User
	db.Find(&users)
	return c.JSON(fiber.Map{"status": "success", "message": "All users", "data": users})
}

// GetUser get a user
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var user model.User
	db.Find(&user, id)
	if user.Username == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "data": user})
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) error {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	db := database.DB
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "errors": err.Error()})
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body", "errors": err.Error()})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "errors": err.Error()})
	}

	user.Password = hash
	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "errors": err.Error()})
	}

	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

// UpdateUser update user
func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Names    string `json:"names"`
		Password string `json:"password,omitempty"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "errors": err.Error()})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	db := database.DB
	var user model.User

	result := db.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	// Update fields if provided
	if uui.Username != "" {
		user.Username = uui.Username
	}
	if uui.Email != "" {
		user.Email = uui.Email
	}
	if uui.Names != "" {
		user.Names = uui.Names
	}
	if uui.Password != "" {
		hash, err := hashPassword(uui.Password)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "errors": err.Error()})
		}
		user.Password = hash
	}

	if err := db.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't update user", "errors": err.Error()})
	}

	// Return user without password
	type UserResponse struct {
		ID        uint      `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Names     string    `json:"names"`
	}

	userResponse := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Username:  user.Username,
		Email:     user.Email,
		Names:     user.Names,
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": userResponse})
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Review your input", "errors": err.Error()})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	if !validUser(id, pi.Password) {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "Invalid credentials", "data": nil})
	}

	db := database.DB
	var user model.User

	result := db.First(&user, id)
	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't delete user", "errors": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}
