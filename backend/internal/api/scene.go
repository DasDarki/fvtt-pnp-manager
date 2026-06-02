package api

import (
	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type sceneInput struct {
	Name        string         `json:"name"`
	Summary     string         `json:"summary"`
	SceneStatus string         `json:"sceneStatus"`
	SystemData  datatypes.JSON `json:"systemData"`
	FolderID    *uuid.UUID     `json:"folderId"`
}

func (h *Handler) ListScenes(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var rows []models.Scene
	h.db.Where("campaign_id = ?", cam.ID).Order("sort asc, created_at desc").Find(&rows)

	ids := make([]uuid.UUID, 0)
	for _, r := range rows {
		if r.MapAssetID != nil {
			ids = append(ids, *r.MapAssetID)
		}
	}
	if len(ids) > 0 {
		var assets []models.Asset
		h.db.Where("id IN ?", ids).Find(&assets)
		urls := map[uuid.UUID]string{}
		for _, a := range assets {
			urls[a.ID] = a.StorageURL
		}
		for i := range rows {
			if rows[i].MapAssetID != nil {
				rows[i].ImageURL = urls[*rows[i].MapAssetID]
			}
		}
	}
	return c.JSON(rows)
}

func (h *Handler) CreateScene(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in sceneInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name == "" {
		return fail(c, fiber.StatusBadRequest, "name required")
	}
	s := models.Scene{
		CampaignID:  cam.ID,
		Name:        in.Name,
		Summary:     in.Summary,
		SceneStatus: orDefault(in.SceneStatus, "draft"),
		SystemData:  in.SystemData,
		FolderID:    in.FolderID,
	}
	if err := h.db.Create(&s).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not create scene")
	}
	return c.Status(fiber.StatusCreated).JSON(s)
}

func (h *Handler) sceneFor(c *fiber.Ctx) (models.Scene, error) {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return models.Scene{}, err
	}
	var s models.Scene
	err = h.db.Where("id = ? AND campaign_id = ?", c.Params("id"), cam.ID).First(&s).Error
	return s, err
}

func (h *Handler) GetScene(c *fiber.Ctx) error {
	s, err := h.sceneFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "scene not found")
	}
	s.ImageURL = h.imageURLFor(s.MapAssetID)
	return c.JSON(s)
}

func (h *Handler) UpdateScene(c *fiber.Ctx) error {
	s, err := h.sceneFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "scene not found")
	}
	var in sceneInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name != "" {
		s.Name = in.Name
	}
	if in.Summary != "" {
		s.Summary = in.Summary
	}
	if in.SceneStatus != "" {
		s.SceneStatus = in.SceneStatus
	}
	if in.SystemData != nil {
		s.SystemData = in.SystemData
	}
	if in.FolderID != nil {
		s.FolderID = in.FolderID
	}
	s.SyncState = "dirty"
	h.db.Save(&s)
	return c.JSON(s)
}

func (h *Handler) DeleteScene(c *fiber.Ctx) error {
	s, err := h.sceneFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "scene not found")
	}
	h.db.Delete(&s)
	return c.SendStatus(fiber.StatusNoContent)
}
