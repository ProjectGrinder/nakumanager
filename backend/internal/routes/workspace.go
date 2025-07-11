package routes

import (
	"log"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
)

// Interface WorkspaceHandler
type WorkspaceHandler struct {
	Repo     repositories.WorkspaceRepository
	UserRepo repositories.UserRepository
}

// Concrete NewWorkspaceHandler
func NewWorkspaceHandler(workspaceRepo repositories.WorkspaceRepository, userRepo repositories.UserRepository) *WorkspaceHandler {
	return &WorkspaceHandler{
		Repo:     workspaceRepo,
		UserRepo: userRepo,
	}
}

// CreateWorkspace
func (h *WorkspaceHandler) CreateWorkspace(c *fiber.Ctx) error {
	var workspace models.CreateWorkspace

	//Check if user is authenticated
	userID := c.Locals("userID").(string)

	//Parse request
	if err := c.BodyParser(&workspace); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	workspace.ID = uuid.New().String()
	workspace.Name = strings.TrimSpace(workspace.Name)
	workspace.Members = []string{userID}

	if err := validate.Struct(workspace); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errMessages := make([]string, 0, len(validationErrors))
		for _, ve := range validationErrors {
			errMessages = append(errMessages, ve.Field()+" is invalid: "+ve.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMessages})
	}

	//Create workspace
	err := h.Repo.CreateWorkspace(c.Context(), workspace.ID, workspace.Name, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create workspace"})
	}

	//Add creator to workspace members
	err = h.Repo.AddMemberToWorkspace(c.Context(), workspace.ID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add creator to workspace"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "workspace created successfully", "workspace_id": workspace.ID})
}

// GetWorkspacesByUserID
func (h *WorkspaceHandler) GetWorkspacesByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	//Get workspaces
	workspaces, err := h.Repo.ListWorkspacesWithMembersByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch workspaces"})
	}

	return c.Status(fiber.StatusOK).JSON(workspaces)
}

// DeleteWorkspace
func (h *WorkspaceHandler) DeleteWorkspace(c *fiber.Ctx) error {
	//Check if workspace id is provided
	workspaceID := strings.TrimSpace(c.Params("workspaceid"))
	if workspaceID == "" || workspaceID == "undefined" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace id is required"})
	}

	//Check if user is authenticated
	userID := c.Locals("userID").(string)

	//Get workspace
	workspace, err := h.Repo.GetWorkspaceByID(c.Context(), workspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workspace not found"})
	}

	//Check if user is owner
	if workspace.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not authorized to delete this workspace"})
	}

	//Delete workspace
	err = h.Repo.DeleteWorkspace(c.Context(), workspaceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete workspace"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "workspace deleted successfully"})
}

func (h *WorkspaceHandler) UpdateWorkspace(c *fiber.Ctx) error {
	workspaceID := strings.TrimSpace(c.Params("workspaceid"))
	if workspaceID == "" || workspaceID == "undefined" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace id is required"})
	}

	userID := c.Locals("userID").(string)

	workspace, err := h.Repo.GetWorkspaceByID(c.Context(), workspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workspace not found"})
	}

	if workspace.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not authorized to update this workspace"})
	}

	var req models.UpdateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	log.Println("Received update workspace request:", req)

	// Rename
	if req.Name != nil {
		if err := h.Repo.RenameWorkspace(c.Context(), workspaceID, *req.Name); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to rename workspace"})
		}
	}

	// Add members
	if req.AddMembers != nil {
		for _, memberID := range *req.AddMembers {
			if err := h.Repo.AddMemberToWorkspace(c.Context(), workspaceID, memberID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add member to workspace"})
			}
		}
	}

	// Remove members
	if req.RemoveMembers != nil {
		for _, memberID := range *req.RemoveMembers {
			if err := h.Repo.RemoveMemberFromWorkspace(c.Context(), workspaceID, memberID); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove member from workspace"})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "workspace updated successfully"})
}
