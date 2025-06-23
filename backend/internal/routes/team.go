package routes

import (
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
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
func (h *TeamHandler) CreateTeam(c *fiber.Ctx) error {
	var request db.CreateTeamParams
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	workspace, err := h.WorkspaceRepo.GetWorkspaceByID(c.Context(), request.WorkspaceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "workspace not found"})
	}

	if workspace.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission in this workspace"})
	}

	request.ID = uuid.New().String()
	request.Name = strings.TrimSpace(request.Name)

	if err := validate.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errMessages := make([]string, 0, len(validationErrors))
		for _, ve := range validationErrors {
			errMessages = append(errMessages, ve.Field()+" is invalid: "+ve.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMessages})
	}

	err = h.Repo.CreateTeam(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create team"})
	}

	err = h.Repo.AddMemberToTeam(c.Context(), db.AddMemberToTeamParams{
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

func (h *TeamHandler) GetTeamsByUserID(c *fiber.Ctx) error {
	userId, ok := c.Locals("userID").(string)
	if !ok || userId == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	team, err := h.Repo.GetTeamsByUserID(c.Context(), userId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "teams not found for user"})
	}

	if len(team) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no teams found for user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"teams": team,
	})

}

func (h *TeamHandler) DeleteTeam(c *fiber.Ctx) error {
	teamID := c.Params("id")
	if teamID == "" || teamID == "empty" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team ID is required"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	owner, err := h.Repo.GetOwnerByTeamID(c.Context(), teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

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

	err = h.Repo.DeleteTeamFromTeamMembers(c.Context(), teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove team from team members"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "team deleted successfully",
	})
}

func (h *TeamHandler) AddMemberToTeam(c *fiber.Ctx) error {
	var request db.AddMemberToTeamParams
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	owner, err := h.Repo.GetOwnerByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	leader, err := h.Repo.GetLeaderByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	if owner != userID && leader != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission to remove member from this team"})
	}

	exists, err := h.Repo.IsMemberInTeam(c.Context(), request.TeamID, request.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check member"})
	}

	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "member already exists in the team"})
	}

	err = h.Repo.AddMemberToTeam(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add member to team"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member added to team successfully",
	})
}

func (h *TeamHandler) RemoveMemberFromTeam(c *fiber.Ctx) error {
	var request db.RemoveMemberFromTeamParams
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	_, err := h.Repo.GetOwnerByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	owner, err := h.Repo.GetOwnerByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	leader, err := h.Repo.GetLeaderByTeamID(c.Context(), request.TeamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	if owner != userID && leader != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission to remove member from this team"})
	}

	err = h.Repo.RemoveMemberFromTeam(c.Context(), request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove member from team"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member removed from team successfully",
	})
}
