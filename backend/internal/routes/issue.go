package routes

import (
	"sync"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	models "github.com/nack098/nakumanager/internal/models"
)

var (
	Validate = validator.New()
	//mock up DB
	Issues     = make(map[string]models.Issue)
	IssueMutex = sync.RWMutex{}
)

func SetUpIssueRoutes(api fiber.Router) {
	api.Post("/issues", CreateIssue)
	api.Get("/issues/:id", GetIssuesByUserID)
	api.Delete("/issues/:id", DeleteIssue)
}

func CreateIssue(c *fiber.Ctx) error {
	var issue models.Issue

	if err := c.BodyParser(&issue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	ownerID, ok := c.Locals("userID").(string)
	if !ok || ownerID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: missing owner ID"})
	}
	issue.OwnerID = ownerID

	if err := Validate.Struct(issue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	issue.ID = uuid.New().String()

	IssueMutex.Lock()
	Issues[issue.ID] = issue
	IssueMutex.Unlock()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "issue created successfully"})
}

func GetIssuesByUserID(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "UserID is required"})
	}

	var result []models.Issue
	for _, issue := range Issues {
		for _, assignee := range issue.Assignee {
			if assignee == userID {
				result = append(result, issue)
				break
			}
		}
	}

	return c.Status(200).JSON(result)
}

func DeleteIssue(c *fiber.Ctx) error {
	issueID := c.Params("id")

	userIDVal := c.Locals("userID")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	IssueMutex.RLock()
	issue, exists := Issues[issueID]
	IssueMutex.RUnlock()

	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Issue not found"})
	}

	isAssigneeOrOwner := false
	if issue.OwnerID == userID {
		isAssigneeOrOwner = true
	} else {
		for _, a := range issue.Assignee {
			if a == userID {
				isAssigneeOrOwner = true
				break
			}
		}
	}
	if !isAssigneeOrOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	IssueMutex.Lock()
	delete(Issues, issueID)
	defer IssueMutex.Unlock()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Issue deleted successfully"})
}
