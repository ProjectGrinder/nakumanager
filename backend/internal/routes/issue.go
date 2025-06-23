package routes

import (
	"database/sql"
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
)

type IssueHandler struct {
	Repo        repositories.IssueRepository
	TeamRepo    repositories.TeamRepository
	ProjectRepo repositories.ProjectRepository
}

func NewIssueHandler(repo repositories.IssueRepository, teamRepo repositories.TeamRepository, projectRepo repositories.ProjectRepository) *IssueHandler {
	return &IssueHandler{
		Repo:        repo,
		TeamRepo:    teamRepo,
		ProjectRepo: projectRepo,
	}
}

func ToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func (h *IssueHandler) CreateIssue(c *fiber.Ctx) error {
	var issueReq models.IssueCreate
	if err := c.BodyParser(&issueReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	issueReq.OwnerID = userID

	validate := validator.New()
	if err := validate.Struct(&issueReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"detail": err.Error(),
		})
	}

	teamExists, err := h.TeamRepo.IsTeamExists(c.Context(), issueReq.TeamID)
	if err != nil || !teamExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	isMember, err := h.TeamRepo.IsMemberInTeam(c.Context(), issueReq.TeamID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check team membership"})
	}
	if !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not a member of the team"})
	}

	if issueReq.ProjectID != nil {
		projectExists, err := h.ProjectRepo.IsProjectExists(c.Context(), *issueReq.ProjectID)
		if err != nil || !projectExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project not found",
			})
		}
	}

	if issueReq.Assignee != nil {
		userExists, err := h.TeamRepo.IsMemberInTeam(c.Context(), issueReq.TeamID, *issueReq.Assignee)
		if err != nil || !userExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Assignee not found",
			})
		}
	}

	if issueReq.Status == "" {
		issueReq.Status = "todo"
	}

	if issueReq.Priority == nil {
		def := "low"
		issueReq.Priority = &def
	}
	if issueReq.StartDate == nil {
		now := time.Now().UTC()
		issueReq.StartDate = &now
	}

	issueReq.ID = uuid.New().String()
	body := db.CreateIssueParams{
		ID:        issueReq.ID,
		Title:     issueReq.Title,
		Content:   ToNullString(issueReq.Content),
		Priority:  ToNullString(issueReq.Priority),
		Status:    issueReq.Status,
		Assignee:  ToNullString(issueReq.Assignee),
		ProjectID: ToNullString(issueReq.ProjectID),
		TeamID:    issueReq.TeamID,
		StartDate: ToNullTime(issueReq.StartDate),
		EndDate:   ToNullTime(issueReq.EndDate),
		Label:     ToNullString(issueReq.Label),
		OwnerID:   issueReq.OwnerID,
	}

	if err := h.Repo.CreateIssue(c.Context(), body); err != nil {
		log.Printf("Failed to create issue: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create issue",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Issue created successfully",
		"issueID": issueReq.ID,
	})
}

func (h *IssueHandler) AddAssigneeToIssue(c *fiber.Ctx) error {
	var assigneeReq models.AssigneeRequest
	if err := c.BodyParser(&assigneeReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if err := h.Repo.AddAssigneeToIssue(c.Context(), db.AddAssigneeToIssueParams{
		IssueID: assigneeReq.IssueID,
		UserID:  assigneeReq.UserID,
	}); err != nil {
		log.Printf("Failed to add assignee to issue: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add assignee to issue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Assignee added successfully",
	})
}

func (h *IssueHandler) RemoveAssigneeFromIssue(c *fiber.Ctx) error {
	var assigneeReq models.AssigneeRequest
	if err := c.BodyParser(&assigneeReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if err := h.Repo.RemoveAssigneeFromIssue(c.Context(), db.RemoveAssigneeFromIssueParams{
		IssueID: assigneeReq.IssueID,
		UserID:  assigneeReq.UserID,
	}); err != nil {
		log.Printf("Failed to remove assignee from issue: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove assignee from issue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Assignee removed successfully",
	})
}

func (h *IssueHandler) GetIssuesByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}
	issues, err := h.Repo.GetIssueByUserID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get issues by user ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get issues",
		})
	}
	return c.Status(fiber.StatusOK).JSON(issues)
}

func (h *IssueHandler) DeleteIssue(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	issueID := c.Params("id")
	if issueID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Issue ID is required",
		})
	}

	issue, err := h.Repo.GetIssueByID(c.Context(), issueID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Issue not found",
		})
	}

	if issue.OwnerID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this issue",
		})
	}

	if err := h.Repo.DeleteIssue(c.Context(), issueID); err != nil {
		log.Printf("Failed to delete issue: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete issue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Issue deleted successfully",
	})
}
