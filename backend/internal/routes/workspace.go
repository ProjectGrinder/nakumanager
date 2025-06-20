package routes

import (
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
)

var validate = validator.New()

type WorkspaceHandler struct {
	Repo     repositories.WorkspaceRepository
	UserRepo repositories.UserRepository
}

func NewWorkspaceHandler(workspaceRepo repositories.WorkspaceRepository, userRepo repositories.UserRepository) *WorkspaceHandler {
	return &WorkspaceHandler{
		Repo:     workspaceRepo,
		UserRepo: userRepo,
	}
}

func (h *WorkspaceHandler) CreateWorkspace(c *fiber.Ctx) error {
	var workspace models.Workspace

	if err := c.BodyParser(&workspace); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
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

	err := h.Repo.CreateWorkspace(c.Context(), workspace.ID, workspace.Name, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create workspace"})
	}

	err = h.Repo.AddMemberToWorkspace(c.Context(), workspace.ID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add creator to workspace"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "workspace created successfully", "workspace_id": workspace.ID})
}

func (h *WorkspaceHandler) GetWorkspacesByUserID(c *fiber.Ctx) error {
	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	workspaces, err := h.Repo.ListWorkspacesWithMembersByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch workspaces"})
	}

	return c.Status(fiber.StatusOK).JSON(workspaces)
}

func (h *WorkspaceHandler) DeleteWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("workspaceid")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace id is required"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	workspace, err := h.Repo.GetWorkspaceByID(c.Context(), workspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workspace not found"})
	}

	if workspace.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not authorized to delete this workspace"})
	}

	err = h.Repo.DeleteWorkspace(c.Context(), workspaceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete workspace"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "workspace deleted successfully"})
}

func (h *WorkspaceHandler) AddMemberToWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("workspaceid")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace id is required"})
	}

	requesterID, ok := c.Locals("userID").(string)
	if !ok || requesterID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	workspace, err := h.Repo.GetWorkspaceByID(c.Context(), workspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workspace not found"})
	}

	if workspace.OwnerID != requesterID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not authorized to add members"})
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&body); err != nil || body.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	err = h.Repo.AddMemberToWorkspace(c.Context(), workspaceID, body.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add member to workspace"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "member added successfully"})
}

func (h *WorkspaceHandler) RemoveMemberFromWorkspace(c *fiber.Ctx) error {
	workspaceID := c.Params("workspaceid")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace id is required"})
	}

	requesterID, ok := c.Locals("userID").(string)
	if !ok || requesterID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	workspace, err := h.Repo.GetWorkspaceByID(c.Context(), workspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workspace not found"})
	}

	if workspace.OwnerID != requesterID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not authorized to remove members"})
	}

	var body struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&body); err != nil || body.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id is required"})
	}

	err = h.Repo.RemoveMemberFromWorkspace(c.Context(), workspaceID, body.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove member from workspace"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "member removed successfully"})
}

func (h *WorkspaceHandler) RenameWorkSpace(c *fiber.Ctx) error {
	workspaceID := c.Params("workspaceid")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace id is required"})
	}

	requesterID, ok := c.Locals("userID").(string)
	if !ok || requesterID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	workspace, err := h.Repo.GetWorkspaceByID(c.Context(), workspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workspace not found"})
	}

	if workspace.OwnerID != requesterID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not authorized to rename this workspace"})
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&body); err != nil || body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	err = h.Repo.RenameWorkspace(c.Context(), workspaceID, body.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to rename workspace"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "workspace renamed successfully"})
}
