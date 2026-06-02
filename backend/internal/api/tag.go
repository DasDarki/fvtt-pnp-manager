package api

import (
	"strings"

	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type tagInput struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (h *Handler) ListTags(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var rows []models.Tag
	h.db.Where("campaign_id = ?", cam.ID).Order("name asc").Find(&rows)
	return c.JSON(rows)
}

func (h *Handler) CreateTag(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in tagInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	in.Name = strings.TrimSpace(in.Name)
	if in.Name == "" {
		return fail(c, fiber.StatusBadRequest, "name required")
	}
	var existing models.Tag
	if err := h.db.Where("campaign_id = ? AND LOWER(name) = LOWER(?)", cam.ID, in.Name).First(&existing).Error; err == nil {
		return c.JSON(existing)
	}
	tag := models.Tag{CampaignID: cam.ID, Name: in.Name, Color: orDefault(in.Color, "secondary")}
	if err := h.db.Create(&tag).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not create tag")
	}
	return c.Status(fiber.StatusCreated).JSON(tag)
}

func (h *Handler) DeleteTag(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	tagID := c.Params("id")
	h.db.Where("campaign_id = ? AND tag_id = ?", cam.ID, tagID).Delete(&models.EntityTag{})
	h.db.Where("id = ? AND campaign_id = ?", tagID, cam.ID).Delete(&models.Tag{})
	return c.SendStatus(fiber.StatusNoContent)
}

type entityTagInput struct {
	TagID      uuid.UUID `json:"tagId"`
	EntityType string    `json:"entityType"`
	EntityID   uuid.UUID `json:"entityId"`
}

type entityTagOut struct {
	ID         uuid.UUID  `json:"id"`
	EntityID   uuid.UUID  `json:"entityId"`
	EntityType string     `json:"entityType"`
	Tag        models.Tag `json:"tag"`
}

func (h *Handler) ListEntityTags(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	q := h.db.Where("campaign_id = ?", cam.ID)
	if subject := c.Query("subjectId"); subject != "" {
		q = q.Where("entity_id = ?", subject)
	}
	if et := c.Query("entityType"); et != "" {
		q = q.Where("entity_type = ?", et)
	}
	var rows []models.EntityTag
	q.Find(&rows)

	tagIDs := make([]uuid.UUID, 0, len(rows))
	for _, r := range rows {
		tagIDs = append(tagIDs, r.TagID)
	}
	tagMap := map[uuid.UUID]models.Tag{}
	if len(tagIDs) > 0 {
		var tags []models.Tag
		h.db.Where("id IN ?", tagIDs).Find(&tags)
		for _, t := range tags {
			tagMap[t.ID] = t
		}
	}

	out := make([]entityTagOut, 0, len(rows))
	for _, r := range rows {
		t, ok := tagMap[r.TagID]
		if !ok {
			continue
		}
		out = append(out, entityTagOut{ID: r.ID, EntityID: r.EntityID, EntityType: r.EntityType, Tag: t})
	}
	return c.JSON(out)
}

func (h *Handler) CreateEntityTag(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in entityTagInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.TagID == uuid.Nil || in.EntityID == uuid.Nil || !validLinkTypes[in.EntityType] {
		return fail(c, fiber.StatusBadRequest, "valid tagId/entityType/entityId required")
	}
	var existing models.EntityTag
	if err := h.db.Where("campaign_id = ? AND tag_id = ? AND entity_type = ? AND entity_id = ?",
		cam.ID, in.TagID, in.EntityType, in.EntityID).First(&existing).Error; err == nil {
		return c.JSON(existing)
	}
	et := models.EntityTag{CampaignID: cam.ID, TagID: in.TagID, EntityType: in.EntityType, EntityID: in.EntityID}
	if err := h.db.Create(&et).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not attach tag")
	}
	return c.Status(fiber.StatusCreated).JSON(et)
}

func (h *Handler) DeleteEntityTag(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	h.db.Where("id = ? AND campaign_id = ?", c.Params("id"), cam.ID).Delete(&models.EntityTag{})
	return c.SendStatus(fiber.StatusNoContent)
}
