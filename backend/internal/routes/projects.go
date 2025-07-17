package routes

import (
	"database/sql"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
	"github.com/nack098/nakumanager/internal/ws"
)

type ProjectHandler struct {
	DB       *sql.DB
	Repo     repositories.ProjectRepository
	TeamRepo repositories.TeamRepository
}

func NewProjectHandler(db *sql.DB, repo repositories.ProjectRepository, teamRepo repositories.TeamRepository) *ProjectHandler {
	return &ProjectHandler{
		DB:       db,
		Repo:     repo,
		TeamRepo: teamRepo,
	}
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var body models.CreateProject

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Check if the user is authenticated
	userID := c.Locals("userID").(string)

	// ตรวจสอบทีมมีอยู่จริงและอยู่ใน workspace เดียวกัน
	team, err := h.TeamRepo.GetTeamByID(c.Context(), body.TeamID)
	if err != nil {
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
	body.ID = projectID
	if body.LeaderID != nil && strings.TrimSpace(*body.LeaderID) == "" {
		body.LeaderID = nil
	}
	body.CreatedBy = userID

	// สร้างโปรเจกต์
	if err := h.Repo.CreateProject(c.Context(), body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Project created successfully"})
}

func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	projectID := c.Params("id")
	if projectID == "" || projectID == "undefined" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing project ID"})
	}

	var body models.EditProject
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID := c.Locals("userID").(string)

	// ตรวจสอบงาสคนที่แก้ไขคือ Owner ไม่ก็ Leader
	owner, err := h.Repo.GetOwnerByProjectID(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
	}
	if owner != userID {
		leader, err := h.Repo.GetLeaderByProjectID(c.Context(), projectID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
		}
		if leader != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "you are not authorized to update this project"})
		}
	}

	// ตรวจสอบว่า project มีอยู่จริง
	_, err = h.Repo.GetProjectByID(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch project"})
	}

	body.ID = projectID
	query, args := buildUpdateQuery(body)

	if body.AddMember != nil && len(*body.AddMember) > 0 {
		for _, member := range *body.AddMember {
			h.Repo.AddMemberToProject(c.Context(), projectID, member)
		}
	}

	if body.RemoveMember != nil && len(*body.RemoveMember) > 0 {
		for _, member := range *body.RemoveMember {
			h.Repo.RemoveMemberFromProject(c.Context(), projectID, member)
		}
	}

	if _, err := h.DB.ExecContext(c.Context(), query, args...); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update project"})
	}

	ws.BroadcastToRoom("project", projectID, "project_updated", body)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Project updated successfully"})
}

func (h *ProjectHandler) GetProjectsByUserID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	// ดึงโปรเจกต์
	projects, err := h.Repo.GetProjectsByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch projects"})
	}

	return c.JSON(projects)
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	// ตรวจสอบ param
	projectID := c.Params("id")
	if strings.TrimSpace(projectID) == "" || projectID == "undefined" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "project_id is required"})
	}

	// ตรวจสอบสิทธิ์
	userID := c.Locals("userID").(string)

	// ตรวจสอบว่า project มีอยู่
	project, err := h.Repo.GetProjectByID(c.Context(), projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch project"})
	}

	// ตรวจสอบว่าผู้ใช้เป็น Owner หรือ Leader
	if project.CreatedBy != userID && project.LeaderID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only project creator or leader can delete project"})
	}

	// ลบ
	if err := h.Repo.DeleteProject(c.Context(), projectID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete project"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Project deleted successfully"})
}
