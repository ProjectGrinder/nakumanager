package routes

import (
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
)

type TeamHandler struct {
	Repo          repositories.TeamRepository
	WorkspaceRepo repositories.WorkspaceRepository
}

func NewTeamHandler(repo repositories.TeamRepository, workspaceRepo repositories.WorkspaceRepository) *TeamHandler {
	return &TeamHandler{
		Repo:          repo,
		WorkspaceRepo: workspaceRepo,
	}
}

// CreateTeam
func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	//Check if user is authenticated
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	//Parse request
	var request models.CreateTeam
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	//Check if workspace exists
	workspace, err := h.WorkspaceRepo.GetWorkspaceByID(c.Context(), request.WorkspaceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace not found"})
	}

	//Check if user is owner of workspace
	if workspace.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission in this workspace"})
	}

	request.ID = uuid.New().String()
	request.Name = strings.TrimSpace(request.Name)

	//Validate
	if err := validate.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errMessages := make([]string, 0, len(validationErrors))
		for _, ve := range validationErrors {
			errMessages = append(errMessages, ve.Field()+" is invalid: "+ve.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMessages})
	}

	//Create team
	err = h.Repo.CreateTeam(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create team"})
	}

	//Add user to team
	err = h.Repo.AddMemberToTeam(c.Context(), models.AddMemberToTeam{
		TeamID: request.ID,
		UserID: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add member to team"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "team created successfully",
	})
}

// GetTeamsByUserID
func (h *TeamHandler) GetTeamsByUserID(c *fiber.Ctx) error {
	//Check if user is authenticated
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	//Get teams
	team, err := h.Repo.GetTeamsByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "teams not found for user"})
	}

	//Check if teams exist
	if len(team) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no teams found for user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"teams": team,
	})

}

// DeleteTeam
func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	//Check if user is authenticated
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	//Parse request
	teamID := c.Params("id")
	if teamID == "" || teamID == "empty" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team ID is required"})
	}

	//Check if user is owner
	owner, err := h.Repo.GetOwnerByTeamID(c.Context(), teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	//Check if user is leader
	leader, err := h.Repo.GetLeaderByTeamID(c.Context(), teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	if owner != userID && leader != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission to remove member from this team"})
	}

	err = h.Repo.DeleteTeam(c.Context(), teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete team"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "team deleted successfully",
	})
}

func (h *TeamHandler) AddMemberToTeam(c *fiber.Ctx) error {
	//Check if user is authenticated
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	teamID := c.Params("id")
	if teamID == "" || teamID == "empty" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team ID is required"})
	}

	//Parse request
	var request models.AddMemberToTeam
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	request.TeamID = teamID

	owner, err := h.Repo.GetOwnerByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	leader, err := h.Repo.GetLeaderByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	//Check if user is not owner and leader
	if owner != userID && leader != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission to remove member from this team"})
	}

	//Check if member already exists
	exists, err := h.Repo.IsMemberInTeam(c.Context(), request.TeamID, request.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check member"})
	}

	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "member already exists in the team"})
	}

	//Add member
	err = h.Repo.AddMemberToTeam(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add member to team"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member added to team successfully",
	})
}

func (h *TeamHandler) RemoveMemberFromTeam(c *fiber.Ctx) error {
	//Check if user is authenticated
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	teamID := c.Params("id")
	if teamID == "" || teamID == "empty" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team ID is required"})
	}

	//Parse request
	var request models.RemoveMemberFromTeam
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	request.TeamID = teamID

	owner, err := h.Repo.GetOwnerByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	leader, err := h.Repo.GetLeaderByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	//Check if user is not owner and leader
	if owner != userID && leader != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission to remove member from this team"})
	}

	//Remove member
	err = h.Repo.RemoveMemberFromTeam(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove member from team"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member removed from team successfully",
	})
}

func (h *TeamHandler) RenameTeam(c *fiber.Ctx) error {
	//Check if user is authenticated
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	teamID := c.Params("id")
	if teamID == "" || teamID == "empty" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team ID is required"})
	}

	//Parse request
	var request models.RenameTeam
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	request.TeamID = teamID

	//Check if team exists
	exists, err := h.Repo.IsTeamExists(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check team"})
	}

	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	//Check if user is owner
	owner, err := h.Repo.GetOwnerByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	leader, err := h.Repo.GetLeaderByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	//Check if user is not owner and leader
	if owner != userID && leader != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission to remove member from this team"})
	}

	//Rename team
	err = h.Repo.RenameTeam(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to rename team"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "team renamed successfully",
	})
}

func (h *TeamHandler) SetTeamLeader(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	teamID := c.Params("id")
	if teamID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team ID is required"})
	}

	var request models.SetTeamLeader
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	request.TeamID = teamID

	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "validation failed"})
	}

	ctx := c.Context()

	exists, err := h.Repo.IsTeamExists(ctx, request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check team"})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	owner, err := h.Repo.GetOwnerByTeamID(ctx, request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}
	if owner != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission to set team leader"})
	}

	if err := h.Repo.SetLeaderToTeam(ctx, request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to set team leader"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "team leader set successfully",
	})
}
