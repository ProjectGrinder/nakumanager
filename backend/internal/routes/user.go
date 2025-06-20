package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/repositories"
)

type UserManage struct {
	repo repositories.UserRepository
}

func NewUserManage(repo repositories.UserRepository) *UserManage {
	return &UserManage{repo: repo}
}

func (h *UserManage) CreateUser(c *fiber.Ctx) error {
	var req models.User
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	id := uuid.New().String()

	arg := db.CreateUserParams{
		ID:           id,
		Username:     req.Username,
		PasswordHash: req.PasswordHash,
		Email:        req.Email,
		Roles:        req.Roles,
	}

	user := h.repo.CreateUser(c.Context(), arg)
	if user == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create user")
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserManage) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.repo.ListUsers(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(users)
}

func (h *UserManage) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.repo.GetUserByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	return c.JSON(user)
}

func (h *UserManage) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.repo.DeleteUser(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
