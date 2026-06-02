package api

import (
	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validLinkTypes = map[string]bool{"character": true, "item": true, "scene": true, "image": true}

type linkInput struct {
	FromType string    `json:"fromType"`
	FromID   uuid.UUID `json:"fromId"`
	ToType   string    `json:"toType"`
	ToID     uuid.UUID `json:"toId"`
	Kind     string    `json:"kind"`
}

type linkRef struct {
	Type  string    `json:"type"`
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Image string    `json:"image"`
}

type linkOut struct {
	ID    uuid.UUID `json:"id"`
	Kind  string    `json:"kind"`
	Other linkRef   `json:"other"`
}

type refKey struct {
	Type string
	ID   uuid.UUID
}

func (h *Handler) resolveRefs(campaignID uuid.UUID, refs []refKey) map[string]linkRef {
	byType := map[string][]uuid.UUID{}
	for _, r := range refs {
		byType[r.Type] = append(byType[r.Type], r.ID)
	}
	out := map[string]linkRef{}
	key := func(t string, id uuid.UUID) string { return t + ":" + id.String() }

	for typ, ids := range byType {
		switch typ {
		case "character":
			var rows []models.Character
			h.db.Where("campaign_id = ? AND id IN ?", campaignID, ids).Find(&rows)
			for _, x := range rows {
				out[key("character", x.ID)] = linkRef{"character", x.ID, x.Name, h.imageURLFor(x.ImageAssetID)}
			}
		case "item":
			var rows []models.Item
			h.db.Where("campaign_id = ? AND id IN ?", campaignID, ids).Find(&rows)
			for _, x := range rows {
				out[key("item", x.ID)] = linkRef{"item", x.ID, x.Name, h.imageURLFor(x.ImageAssetID)}
			}
		case "scene":
			var rows []models.Scene
			h.db.Where("campaign_id = ? AND id IN ?", campaignID, ids).Find(&rows)
			for _, x := range rows {
				out[key("scene", x.ID)] = linkRef{"scene", x.ID, x.Name, h.imageURLFor(x.MapAssetID)}
			}
		case "image":
			var rows []models.ImageEntry
			h.db.Where("campaign_id = ? AND id IN ?", campaignID, ids).Find(&rows)
			for _, x := range rows {
				out[key("image", x.ID)] = linkRef{"image", x.ID, x.Name, h.imageURLFor(x.AssetID)}
			}
		}
	}
	return out
}

func (h *Handler) ListLinks(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	subject := c.Query("subjectId")
	if subject == "" {
		return fail(c, fiber.StatusBadRequest, "subjectId required")
	}
	sid, err := uuid.Parse(subject)
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid subjectId")
	}

	var links []models.Link
	h.db.Where("campaign_id = ? AND (from_id = ? OR to_id = ?)", cam.ID, sid, sid).
		Order("created_at desc").Find(&links)

	refs := make([]refKey, 0, len(links))
	for _, l := range links {
		if l.FromID == sid {
			refs = append(refs, refKey{l.ToType, l.ToID})
		} else {
			refs = append(refs, refKey{l.FromType, l.FromID})
		}
	}
	resolved := h.resolveRefs(cam.ID, refs)

	out := make([]linkOut, 0, len(links))
	for _, l := range links {
		var r refKey
		if l.FromID == sid {
			r = refKey{l.ToType, l.ToID}
		} else {
			r = refKey{l.FromType, l.FromID}
		}
		ref, ok := resolved[r.Type+":"+r.ID.String()]
		if !ok {
			continue
		}
		out = append(out, linkOut{ID: l.ID, Kind: l.Kind, Other: ref})
	}
	return c.JSON(out)
}

func (h *Handler) CreateLink(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in linkInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if !validLinkTypes[in.FromType] || !validLinkTypes[in.ToType] || in.FromID == uuid.Nil || in.ToID == uuid.Nil {
		return fail(c, fiber.StatusBadRequest, "valid fromType/fromId/toType/toId required")
	}
	if in.FromType == in.ToType && in.FromID == in.ToID {
		return fail(c, fiber.StatusBadRequest, "cannot link an entity to itself")
	}

	var existing models.Link
	err = h.db.Where(
		"campaign_id = ? AND ((from_type = ? AND from_id = ? AND to_type = ? AND to_id = ?) OR (from_type = ? AND from_id = ? AND to_type = ? AND to_id = ?))",
		cam.ID,
		in.FromType, in.FromID, in.ToType, in.ToID,
		in.ToType, in.ToID, in.FromType, in.FromID,
	).First(&existing).Error
	if err == nil {
		return c.JSON(existing)
	}

	link := models.Link{
		CampaignID: cam.ID,
		FromType:   in.FromType,
		FromID:     in.FromID,
		ToType:     in.ToType,
		ToID:       in.ToID,
		Kind:       orDefault(in.Kind, "related"),
	}
	if err := h.db.Create(&link).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not create link")
	}
	return c.Status(fiber.StatusCreated).JSON(link)
}

func (h *Handler) DeleteLink(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	h.db.Where("id = ? AND campaign_id = ?", c.Params("id"), cam.ID).Delete(&models.Link{})
	return c.SendStatus(fiber.StatusNoContent)
}
