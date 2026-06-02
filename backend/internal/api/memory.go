package api

import (
	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type memoryInput struct {
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Kind        string     `json:"kind"`
	Level        string     `json:"level"`
	SubjectType  string     `json:"subjectType"`
	SubjectID    *uuid.UUID `json:"subjectId"`
	SubjectLabel string     `json:"subjectLabel"`
	Pinned       bool       `json:"pinned"`
}

type memoryPatch struct {
	Title        *string `json:"title"`
	Body         *string `json:"body"`
	Level        *string `json:"level"`
	Acknowledged *bool   `json:"acknowledged"`
	Pinned       *bool   `json:"pinned"`
}

var validLevels = map[string]bool{"info": true, "notice": true, "warning": true, "critical": true}

func (h *Handler) ListMemories(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	q := h.db.Where("campaign_id = ?", cam.ID)
	if subject := c.Query("subjectId"); subject != "" {
		q = q.Where("subject_id = ?", subject)
	}
	var rows []models.Memory
	q.Order("pinned desc, created_at desc").Find(&rows)
	return c.JSON(rows)
}

func (h *Handler) ListAlerts(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var rows []models.Memory
	h.db.Where("campaign_id = ? AND level = ? AND acknowledged = ?", cam.ID, "critical", false).
		Order("created_at desc").Find(&rows)
	return c.JSON(rows)
}

func (h *Handler) CreateMemory(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in memoryInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Title == "" {
		return fail(c, fiber.StatusBadRequest, "title required")
	}
	level := orDefault(in.Level, "info")
	if !validLevels[level] {
		return fail(c, fiber.StatusBadRequest, "invalid level")
	}
	m := models.Memory{
		CampaignID:  cam.ID,
		Title:       in.Title,
		Body:        in.Body,
		Kind:        orDefault(in.Kind, "note"),
		Level:        level,
		SubjectType:  orDefault(in.SubjectType, "campaign"),
		SubjectID:    in.SubjectID,
		SubjectLabel: in.SubjectLabel,
		Pinned:       in.Pinned,
	}
	if err := h.db.Create(&m).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not create memory")
	}
	return c.Status(fiber.StatusCreated).JSON(m)
}

func (h *Handler) memoryFor(c *fiber.Ctx) (models.Memory, error) {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return models.Memory{}, err
	}
	var m models.Memory
	err = h.db.Where("id = ? AND campaign_id = ?", c.Params("id"), cam.ID).First(&m).Error
	return m, err
}

func (h *Handler) UpdateMemory(c *fiber.Ctx) error {
	m, err := h.memoryFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "memory not found")
	}
	var in memoryPatch
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Title != nil {
		m.Title = *in.Title
	}
	if in.Body != nil {
		m.Body = *in.Body
	}
	if in.Level != nil {
		if !validLevels[*in.Level] {
			return fail(c, fiber.StatusBadRequest, "invalid level")
		}
		m.Level = *in.Level
	}
	if in.Acknowledged != nil {
		m.Acknowledged = *in.Acknowledged
	}
	if in.Pinned != nil {
		m.Pinned = *in.Pinned
	}
	h.db.Save(&m)
	return c.JSON(m)
}

func (h *Handler) DeleteMemory(c *fiber.Ctx) error {
	m, err := h.memoryFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "memory not found")
	}
	h.db.Delete(&m)
	return c.SendStatus(fiber.StatusNoContent)
}
