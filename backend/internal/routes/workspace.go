package routes

import (
	"strings"
	"sync"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	models "github.com/nack098/nakumanager/internal/models"
)

var (
	WorkSpaces     = make(map[string]models.Workspace)
	WorkspaceMutex = sync.RWMutex{}
	validate       = validator.New()
)

func SetUpWorkspaceRoutes(api fiber.Router) {
	api.Get("/workspace/:id", GetWorkspacesByUserID)
	api.Post("/workspace", CreateWorkspace)
	api.Delete("/workspace/:workspaceid", DeleteWorkspace)
}

func CreateWorkspace(c *fiber.Ctx) error {
	var workspace models.Workspace

	if err := c.BodyParser(&workspace); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	workspace.ID = uuid.New().String()

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	workspace.Members = []string{userID}

	workspace.Name = strings.TrimSpace(workspace.Name)

	if err := validate.Struct(workspace); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errMessages := make([]string, 0, len(validationErrors))
		for _, ve := range validationErrors {
			errMessages = append(errMessages, ve.Field()+" is invalid: "+ve.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMessages})
	}

	WorkspaceMutex.Lock()
	WorkSpaces[workspace.ID] = workspace
	WorkspaceMutex.Unlock()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "workspace created successfully"})
}

func GetWorkspacesByUserID(c *fiber.Ctx) error {
	userId := c.Params("id")
	if userId == "" {
		return c.Status(400).JSON(fiber.Map{"error": "UserID is required"})
	}

	var result []models.Workspace
	for _, workspace := range WorkSpaces {
		for _, member := range workspace.Members {
			if member == userId {
				result = append(result, workspace)
				break
			}
		}
	}

	return c.Status(200).JSON(result)
}
func DeleteWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("workspaceid")

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	WorkspaceMutex.RLock()
	workspace, exists := WorkSpaces[workspaceID]
	WorkspaceMutex.RUnlock()

	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Workspace not found"})
	}

	user := models.GetUserByID(userID)
	if user.ID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}

	if user.Roles != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to delete this workspace"})
	}

	WorkspaceMutex.Lock()
	delete(WorkSpaces, workspace.ID)
	WorkspaceMutex.Unlock()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Workspace deleted successfully",
	})
}
