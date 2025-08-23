package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"app/database"
	"app/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// ---------------------------------------------------------------------
//  ENDPOINTS (mounted in router/router.go)
//  GET  /api/resource           – list all resources
//  GET  /api/resource/:id       – get one resource
//  POST /api/resource           – create a new resource (JWT protected)
//  PUT  /api/resource/:id       – update resource (JWT protected)
//  DELETE /api/resource/:id     – delete resource (JWT protected)
//  GET  /api/resource/:id/history – get resource change history
// ---------------------------------------------------------------------

// resourceCreateInput describes the JSON payload for creating resources
type resourceCreateInput struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description"`
	Unit        string `json:"unit" validate:"required,min=1,max=20"`
	Quantity    int    `json:"quantity" validate:"min=0"`
}

// resourceUpdateInput describes the JSON payload for updating resources
type resourceUpdateInput struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty"`
	Unit        *string `json:"unit,omitempty" validate:"omitempty,min=1,max=20"`
	Quantity    *int    `json:"quantity,omitempty" validate:"omitempty,min=0"`
}

// logResourceChange logs changes to the resource history table
func logResourceChange(resourceID uint, action string, userID uint, oldData, newData interface{}, description string) error {
	db := database.DB

	history := model.ResourceHistory{
		ResourceID:  resourceID,
		Action:      action,
		UserID:      userID,
		Timestamp:   time.Now(),
		Description: description,
	}

	if oldData != nil {
		oldJSON, err := json.Marshal(oldData)
		if err != nil {
			return err
		}
		history.OldData = string(oldJSON)
	}

	if newData != nil {
		newJSON, err := json.Marshal(newData)
		if err != nil {
			return err
		}
		history.NewData = string(newJSON)
	}

	return db.Create(&history).Error
}

// getUserIDFromToken extracts user ID from JWT token
func getUserIDFromToken(c *fiber.Ctx) uint {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0
	}
	return uint(userIDFloat)
}

// ----------  GET ALL --------------------------------------------------

// GetAllResources returns a list of all resources
func GetAllResources(c *fiber.Ctx) error {
	db := database.DB
	var resources []model.Resource

	if err := db.Find(&resources).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "cannot fetch resources", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "resources list", "data": resources})
}

// ----------  GET ONE --------------------------------------------------

// GetResource returns a single resource by its numeric ID
func GetResource(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var resource model.Resource

	if err := db.First(&resource, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"status": "error", "message": "resource not found", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "resource found", "data": resource})
}

// ----------  CREATE ---------------------------------------------------

// CreateResource creates a new resource and logs the action
func CreateResource(c *fiber.Ctx) error {
	var input resourceCreateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "invalid json payload", "data": err.Error()})
	}

	if err := validator.New().Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "validation failed", "data": err.Error()})
	}

	// Create the resource
	resource := model.Resource{
		Name:        input.Name,
		Description: input.Description,
		Unit:        input.Unit,
		Quantity:    input.Quantity,
	}

	db := database.DB

	// Begin transaction
	tx := db.Begin()

	// Log to history first
	userID := getUserIDFromToken(c)
	if err := tx.Create(&resource).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "cannot create resource", "data": err.Error()})
	}

	// Log the change
	history := model.ResourceHistory{
		ResourceID:  resource.ID,
		Action:      "CREATE",
		UserID:      userID,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Resource '%s' created", resource.Name),
	}

	newJSON, _ := json.Marshal(resource)
	history.NewData = string(newJSON)

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "cannot log resource change", "data": err.Error()})
	}

	// Commit transaction
	tx.Commit()

	return c.JSON(fiber.Map{"status": "success", "message": "resource created", "data": resource})
}

// ----------  UPDATE ---------------------------------------------------

// UpdateResource partially updates a resource and logs the action
func UpdateResource(c *fiber.Ctx) error {
	id := c.Params("id")
	var input resourceUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "invalid json payload", "data": err.Error()})
	}

	if err := validator.New().Struct(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "validation failed", "data": err.Error()})
	}

	db := database.DB
	var resource model.Resource

	// Get the current resource data
	if err := db.First(&resource, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"status": "error", "message": "resource not found", "data": nil})
	}

	// Store old data for history
	oldResource := resource

	// Apply updates
	if input.Name != nil {
		resource.Name = *input.Name
	}
	if input.Description != nil {
		resource.Description = *input.Description
	}
	if input.Unit != nil {
		resource.Unit = *input.Unit
	}
	if input.Quantity != nil {
		resource.Quantity = *input.Quantity
	}

	// Begin transaction
	tx := db.Begin()

	// Log to history first
	userID := getUserIDFromToken(c)
	history := model.ResourceHistory{
		ResourceID:  resource.ID,
		Action:      "UPDATE",
		UserID:      userID,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Resource '%s' updated", resource.Name),
	}

	oldJSON, _ := json.Marshal(oldResource)
	newJSON, _ := json.Marshal(resource)
	history.OldData = string(oldJSON)
	history.NewData = string(newJSON)

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "cannot log resource change", "data": err.Error()})
	}

	// Update the resource
	if err := tx.Save(&resource).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "cannot update resource", "data": err.Error()})
	}

	// Commit transaction
	tx.Commit()

	return c.JSON(fiber.Map{"status": "success", "message": "resource updated", "data": resource})
}

// ----------  DELETE ---------------------------------------------------

// DeleteResource removes a resource and logs the action
func DeleteResource(c *fiber.Ctx) error {
	id := c.Params("id")

	// Basic sanity check – id must be numeric
	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "invalid resource id", "data": err.Error()})
	}

	db := database.DB
	var resource model.Resource

	// Get the resource data before deletion
	if err := db.First(&resource, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"status": "error", "message": "resource not found", "data": nil})
	}

	// Begin transaction
	tx := db.Begin()

	// Log to history first
	userID := getUserIDFromToken(c)
	history := model.ResourceHistory{
		ResourceID:  resource.ID,
		Action:      "DELETE",
		UserID:      userID,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Resource '%s' deleted", resource.Name),
	}

	oldJSON, _ := json.Marshal(resource)
	history.OldData = string(oldJSON)

	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "cannot log resource change", "data": err.Error()})
	}

	// Delete the resource
	if err := tx.Delete(&resource).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "cannot delete resource", "data": err.Error()})
	}

	// Commit transaction
	tx.Commit()

	return c.JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("resource %s deleted", id), "data": nil})
}

// ----------  HISTORY --------------------------------------------------

// GetResourceHistory returns the change history for a specific resource
func GetResourceHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// Check if resource exists
	var resource model.Resource
	if err := db.First(&resource, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"status": "error", "message": "resource not found", "data": nil})
	}

	// Get history records for this resource
	var history []model.ResourceHistory
	if err := db.Preload("User").Where("resource_id = ?", id).Order("timestamp desc").Find(&history).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "cannot fetch resource history", "data": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "resource history", "data": history})
}
