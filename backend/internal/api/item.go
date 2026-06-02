package api

import (
	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type itemInput struct {
	Name       string         `json:"name"`
	Summary    string         `json:"summary"`
	ItemType   string         `json:"itemType"`
	Rarity     string         `json:"rarity"`
	Attuned    bool           `json:"attuned"`
	SystemData datatypes.JSON `json:"systemData"`
	FolderID   *uuid.UUID     `json:"folderId"`
}

func (h *Handler) ListItems(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var rows []models.Item
	h.db.Where("campaign_id = ?", cam.ID).Order("sort asc, created_at desc").Find(&rows)

	ids := make([]uuid.UUID, 0)
	for _, r := range rows {
		if r.ImageAssetID != nil {
			ids = append(ids, *r.ImageAssetID)
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
			if rows[i].ImageAssetID != nil {
				rows[i].ImageURL = urls[*rows[i].ImageAssetID]
			}
		}
	}
	return c.JSON(rows)
}

func (h *Handler) CreateItem(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in itemInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name == "" {
		return fail(c, fiber.StatusBadRequest, "name required")
	}
	it := models.Item{
		CampaignID: cam.ID,
		Name:       in.Name,
		Summary:    in.Summary,
		ItemType:   in.ItemType,
		Rarity:     orDefault(in.Rarity, "common"),
		Attuned:    in.Attuned,
		SystemData: in.SystemData,
		FolderID:   in.FolderID,
	}
	if err := h.db.Create(&it).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not create item")
	}
	return c.Status(fiber.StatusCreated).JSON(it)
}

func (h *Handler) itemFor(c *fiber.Ctx) (models.Item, error) {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return models.Item{}, err
	}
	var it models.Item
	err = h.db.Where("id = ? AND campaign_id = ?", c.Params("id"), cam.ID).First(&it).Error
	return it, err
}

func (h *Handler) GetItem(c *fiber.Ctx) error {
	it, err := h.itemFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "item not found")
	}
	it.ImageURL = h.imageURLFor(it.ImageAssetID)
	return c.JSON(it)
}

func (h *Handler) UpdateItem(c *fiber.Ctx) error {
	it, err := h.itemFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "item not found")
	}
	var in itemInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name != "" {
		it.Name = in.Name
	}
	if in.Summary != "" {
		it.Summary = in.Summary
	}
	if in.ItemType != "" {
		it.ItemType = in.ItemType
	}
	if in.Rarity != "" {
		it.Rarity = in.Rarity
	}
	if in.SystemData != nil {
		it.SystemData = in.SystemData
	}
	if in.FolderID != nil {
		it.FolderID = in.FolderID
	}
	it.Attuned = in.Attuned
	it.SyncState = "dirty"
	h.db.Save(&it)
	return c.JSON(it)
}

func (h *Handler) DeleteItem(c *fiber.Ctx) error {
	it, err := h.itemFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "item not found")
	}
	h.db.Delete(&it)
	return c.SendStatus(fiber.StatusNoContent)
}
