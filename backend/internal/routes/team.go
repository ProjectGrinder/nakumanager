package routes

import (
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
	"github.com/nack098/nakumanager/internal/ws"
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
	userID := c.Locals("userID").(string)
	//Parse request
	var request models.CreateTeam
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	//Check if workspace exists
	workspace, err := h.WorkspaceRepo.GetWorkspaceByID(c.Context(), request.WorkspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "workspace not found"})
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

// GetTeamsByUserID
func (h *TeamHandler) GetTeamsByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
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
	userID := c.Locals("userID").(string)

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

func (h *TeamHandler) UpdateTeam(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	teamID := strings.TrimSpace(c.Params("id"))
	if teamID == "" || teamID == "empty" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team ID is required"})
	}

	var req models.UpdateTeamRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	ctx := c.Context()

	// Check team existence
	exists, err := h.Repo.IsTeamExists(ctx, teamID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check team"})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	// Get owner and leader
	owner, err := h.Repo.GetOwnerByTeamID(ctx, teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}
	leader, err := h.Repo.GetLeaderByTeamID(ctx, teamID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	}

	if owner != userID && leader != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "no permission to remove member from this team"})
	}

	// Rename team
	if req.Name != nil {
		if err := h.Repo.RenameTeam(ctx, db.RenameTeamParams{ID: teamID, Name: *req.Name}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to rename team"})
		}
	}

	// Add members
	if req.AddMembers != nil {
		for _, memberID := range *req.AddMembers {
			exists, err := h.Repo.IsMemberInTeam(ctx, teamID, memberID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check member"})
			}
			if exists {
				continue
			}
			if err := h.Repo.AddMemberToTeam(ctx, db.AddMemberToTeamParams{TeamID: teamID, UserID: memberID}); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add member"})
			}
		}
	}

	// Remove members
	if req.RemoveMembers != nil {
		for _, memberID := range *req.RemoveMembers {
			if err := h.Repo.RemoveMemberFromTeam(ctx, db.RemoveMemberFromTeamParams{TeamID: teamID, UserID: memberID}); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove member"})
			}
		}
	}

	// Set new leader
	if req.NewLeaderID != nil {
		exists, err := h.Repo.IsMemberInTeam(ctx, teamID, *req.NewLeaderID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check if user is member"})
		}
		if !exists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "new leader must be a member of the team"})
		}
		if err := h.Repo.SetLeaderToTeam(ctx, db.SetLeaderToTeamParams{ID: teamID, LeaderID: *req.NewLeaderID}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to set team leader"})
		}
	}

	ws.BroadcastToRoom("team", teamID, "team_updated", req)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "team updated successfully"})
}
