package routes

import (
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
)

type ProjectHandler struct {
	Repo     repositories.ProjectRepository
	TeamRepo repositories.TeamRepository
}

func NewProjectHandler(repo repositories.ProjectRepository, teamRepo repositories.TeamRepository) *ProjectHandler {
	return &ProjectHandler{
		Repo:     repo,
		TeamRepo: teamRepo,
	}
}

func toNullString(s *string) sql.NullString {
	if s != nil && strings.TrimSpace(*s) != "" {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

func toNullTime(s *string) sql.NullTime {
	if s != nil {
		if t, err := time.Parse(time.RFC3339, *s); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	return sql.NullTime{}
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var body models.CreateProjectRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Validate required fields
	if strings.TrimSpace(body.Name) == "" || strings.TrimSpace(body.WorkspaceID) == "" || strings.TrimSpace(body.TeamID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name, workspace_id, and team_id are required"})
	}

	// ตรวจสอบทีมมีอยู่จริงและอยู่ใน workspace เดียวกัน
	team, err := h.TeamRepo.GetTeamByID(c.Context(), body.TeamID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check team"})
	}

	if team.WorkspaceID != body.WorkspaceID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "team does not belong to the specified workspace"})
	}

	// ตรวจสอบสมาชิกทีม
	isMember, err := h.TeamRepo.IsMemberInTeam(c.Context(), body.TeamID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check team membership"})
	}
	if !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not a member of the team"})
	}

	projectID := uuid.NewString()

	arg := db.CreateProjectParams{
		ID:          projectID,
		Name:        body.Name,
		Status:      toNullString(body.Status),
		Priority:    toNullString(body.Priority),
		WorkspaceID: body.WorkspaceID,
		TeamID:      body.TeamID,
		LeaderID:    toNullString(body.LeaderID),
		StartDate:   toNullTime(body.StartDate),
		EndDate:     toNullTime(body.EndDate),
		Label:       toNullString(body.Label),
		CreatedBy:   userID,
	}

	if err := h.Repo.CreateProject(c.Context(), arg); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create project"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Project created successfully"})
}

func (h *ProjectHandler) AddMemberToProject(c *fiber.Ctx) error {
	var request db.AddMemberToProjectParams
	if err := c.BodyParser(&request); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if request.ProjectID == "" || request.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "project_id and user_id are required"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	project, err := h.Repo.GetProjectByID(c.Context(), request.ProjectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
		}
		log.Println("Error fetching project:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch project"})
	}

	// ตรวจสอบว่า userID เป็น created_by หรือ leader ของโปรเจกต์หรือไม่
	if project.CreatedBy != userID && nullStringToString(project.LeaderID) != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only project creator or leader can add members"})
	}

	// ตรวจสอบว่าผู้ใช้ที่จะเพิ่มอยู่ใน team ของโปรเจกต์หรือไม่
	isMember, err := h.TeamRepo.IsMemberInTeam(c.Context(), project.TeamID, request.UserID)
	if err != nil {
		log.Println("Error checking team membership:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check team membership"})
	}
	if !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user is not a member of the team"})
	}
	
	// เพิ่มสมาชิก
	if err := h.Repo.AddMemberToProject(c.Context(), request); err != nil {
		log.Println("Error adding member to project:", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "member already exists in project"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to add member to project"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Member added to project successfully"})
}

func (h *ProjectHandler) RemoveProjectMembers(c *fiber.Ctx) error {
	var request db.RemoveMemberFromProjectParams
	if err := c.BodyParser(&request); err != nil {
		log.Println("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if request.ProjectID == "" || request.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "project_id and user_id are required"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	project, err := h.Repo.GetProjectByID(c.Context(), request.ProjectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
		}
		log.Println("Error fetching project:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch project"})
	}

	// ตรวจสอบว่า userID เป็น created_by หรือ leader ของโปรเจกต์หรือไม่
	if project.CreatedBy != userID && nullStringToString(project.LeaderID) != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only project creator or leader can remove members"})
	}

	// ลบสมาชิก
	if err := h.Repo.RemoveMemberFromProject(c.Context(), request); err != nil {
		log.Println("Error removing member from project:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to remove member from project"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Member removed from project successfully"})
}


func (h *ProjectHandler) GetProjectsByUserID(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	projects, err := h.Repo.GetProjectsByUserID(c.Context(), userID)
	if err != nil {
		log.Println("Error fetching projects by user ID:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch projects"})
	}

	return c.JSON(projects)
}


func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	projectID := c.Params("id")
	if strings.TrimSpace(projectID) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "project_id is required"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	project, err := h.Repo.GetProjectByID(c.Context(), projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
		}
		log.Println("Error fetching project:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch project"})
	}

	// ✅ ตรวจสอบสิทธิ์: เฉพาะ creator หรือ leader เท่านั้นที่ลบได้
	if project.CreatedBy != userID && nullStringToString(project.LeaderID) != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only project creator or leader can delete project"})
	}

	// ✅ ลบโปรเจกต์
	if err := h.Repo.DeleteProject(c.Context(), projectID); err != nil {
		log.Println("Error deleting project:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete project"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Project deleted successfully"})
}

