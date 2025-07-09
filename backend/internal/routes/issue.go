package routes

import (
	"database/sql"
	"errors"
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
	DB          *sql.DB
	Repo        repositories.IssueRepository
	TeamRepo    repositories.TeamRepository
	ProjectRepo repositories.ProjectRepository
}

func NewIssueHandler(db *sql.DB, repo repositories.IssueRepository, teamRepo repositories.TeamRepository, projectRepo repositories.ProjectRepository) *IssueHandler {
	return &IssueHandler{
		DB:          db,
		Repo:        repo,
		TeamRepo:    teamRepo,
		ProjectRepo: projectRepo,
	}
}

func (h *IssueHandler) CreateIssue(c *fiber.Ctx) error {
	var issueReq models.IssueCreate
	if err := c.BodyParser(&issueReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	issueReq.OwnerID = userID

	if err := validator.New().Struct(&issueReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":  "Validation failed",
			"detail": err.Error(),
		})
	}

	ctx := c.Context()

	// ตรวจสอบทีม
	teamExists, err := h.TeamRepo.IsTeamExists(ctx, issueReq.TeamID)
	if err != nil || !teamExists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Team not found"})
	}

	isMember, err := h.TeamRepo.IsMemberInTeam(ctx, issueReq.TeamID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check team membership"})
	}
	if !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not a member of the team"})
	}

	// ตรวจสอบโปรเจกต์
	if issueReq.ProjectID != nil {
		projectExists, err := h.ProjectRepo.IsProjectExists(ctx, *issueReq.ProjectID)
		if err != nil || !projectExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project not found"})
		}
	}

	// กำหนดค่า default
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

	// สร้าง issue
	issueReq.ID = uuid.New().String()
	body := db.CreateIssueParams{
		ID:        issueReq.ID,
		Title:     issueReq.Title,
		Content:   ToNullString(issueReq.Content),
		Priority:  ToNullString(issueReq.Priority),
		Status:    issueReq.Status,
		ProjectID: ToNullString(issueReq.ProjectID),
		TeamID:    issueReq.TeamID,
		StartDate: ToNullTime(issueReq.StartDate),
		EndDate:   ToNullTime(issueReq.EndDate),
		Label:     ToNullString(issueReq.Label),
		OwnerID:   issueReq.OwnerID,
	}

	if err := h.Repo.CreateIssue(ctx, body); err != nil {
		log.Printf("Failed to create issue: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create issue"})
	}

	// เพิ่ม assignees ถ้ามี
	if issueReq.Assignee != nil {
		for _, userID := range *issueReq.Assignee {
			isValid, err := h.TeamRepo.IsMemberInTeam(ctx, issueReq.TeamID, userID)
			if err != nil {
				log.Printf("Failed to check assignee %s: %v", userID, err)
				continue
			}
			if !isValid {
				log.Printf("User %s is not a member of the team", userID)
				continue
			}
			err = h.Repo.AddAssigneeToIssue(ctx, db.AddAssigneeToIssueParams{
				IssueID: issueReq.ID,
				UserID:  userID,
			})
			if err != nil {
				log.Printf("Failed to add assignee %s: %v", userID, err)
			}
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Issue created successfully",
		"issueID": issueReq.ID,
	})
}

func (h *IssueHandler) UpdateIssue(c *fiber.Ctx) error {
	var req models.UpdateIssueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	req.ID = c.Params("id")
	if req.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing issue ID",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx := c.Context()

	issue, err := h.Repo.GetIssueByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "issue not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch issue",
		})
	}

	isOwner := issue.OwnerID == userID
	isTeamMember, err := h.TeamRepo.IsMemberInTeam(ctx, issue.TeamID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to check team membership",
		})
	}

	if !isOwner && !isTeamMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you are not authorized to update this issue",
		})
	}

	if req.AddAssignee != nil {
		for _, assigneeID := range *req.AddAssignee {
			valid, err := h.TeamRepo.IsMemberInTeam(ctx, issue.TeamID, assigneeID)
			if err != nil {
				log.Printf("Error checking assignee %s: %v", assigneeID, err)
				continue
			}
			if !valid {
				log.Printf("User %s is not a member of the team %s", assigneeID, issue.TeamID)
				continue
			}
			if err := h.Repo.AddAssigneeToIssue(ctx, db.AddAssigneeToIssueParams{
				IssueID: issue.ID,
				UserID:  assigneeID,
			}); err != nil {
				log.Printf("Error adding assignee %s: %v", assigneeID, err)
			} else {
				log.Printf("Successfully added assignee %s to issue %s", assigneeID, issue.ID)
			}
		}
	}

	if req.RemoveAssignee != nil {
		for _, assigneeID := range *req.RemoveAssignee {
			valid, err := h.TeamRepo.IsMemberInTeam(ctx, issue.TeamID, assigneeID)
			if err != nil {
				log.Printf("Error checking assignee %s: %v", assigneeID, err)
				continue
			}
			if !valid {
				log.Printf("User %s is not a member of the team %s", assigneeID, issue.TeamID)
				continue
			}
			if err := h.Repo.RemoveAssigneeFromIssue(ctx, db.RemoveAssigneeFromIssueParams{
				IssueID: issue.ID,
				UserID:  assigneeID,
			}); err != nil {
				log.Printf("Error removing assignee %s: %v", assigneeID, err)
			} else {
				log.Printf("Successfully removed assignee %s from issue %s", assigneeID, issue.ID)
			}
		}
	}

	query, args := buildUpdateIssueQuery(req)
	if query != "" {
		if _, err := h.DB.ExecContext(ctx, query, args...); err != nil {
			log.Println("Update issue failed:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to update issue",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "issue updated successfully",
	})
}

func (h *IssueHandler) DeleteIssue(c *fiber.Ctx) error {
	issue_id := c.Params("id")
	if issue_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing issue ID",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	ctx := c.Context()

	issue, err := h.Repo.GetIssueByID(ctx, issue_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "issue not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch issue",
		})
	}

	isOwner := issue.OwnerID == userID
	isTeamMember, err := h.TeamRepo.IsMemberInTeam(ctx, issue.TeamID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to check team membership",
		})
	}

	if !isOwner && !isTeamMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "you are not authorized to delete this issue",
		})
	}

	if err := h.Repo.DeleteIssue(ctx, issue_id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete issue",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "issue deleted successfully",
	})
}

func (h *IssueHandler) GetIssuesByUserID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	issues, err := h.Repo.GetIssueByUserID(c.Context(), userID)
	if err != nil {
		log.Printf("Failed to get issues by user ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get issues",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"issues": issues,
	})
}
