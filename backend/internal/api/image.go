package api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type imageInput struct {
	Name     string     `json:"name"`
	Notes    string     `json:"notes"`
	PushAs   string     `json:"pushAs"`
	FolderID *uuid.UUID `json:"folderId"`
}

func (h *Handler) ListImages(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var rows []models.ImageEntry
	h.db.Where("campaign_id = ?", cam.ID).Order("sort asc, created_at desc").Find(&rows)

	ids := make([]uuid.UUID, 0)
	for _, r := range rows {
		if r.AssetID != nil {
			ids = append(ids, *r.AssetID)
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
			if rows[i].AssetID != nil {
				rows[i].ImageURL = urls[*rows[i].AssetID]
			}
		}
	}
	return c.JSON(rows)
}

func (h *Handler) CreateImage(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in imageInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name == "" {
		return fail(c, fiber.StatusBadRequest, "name required")
	}
	im := models.ImageEntry{
		CampaignID: cam.ID,
		Name:       in.Name,
		Notes:      in.Notes,
		PushAs:     orDefault(in.PushAs, "empty_actor"),
		FolderID:   in.FolderID,
	}
	if err := h.db.Create(&im).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not create image")
	}
	return c.Status(fiber.StatusCreated).JSON(im)
}

func (h *Handler) imageFor(c *fiber.Ctx) (models.ImageEntry, error) {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return models.ImageEntry{}, err
	}
	var im models.ImageEntry
	err = h.db.Where("id = ? AND campaign_id = ?", c.Params("id"), cam.ID).First(&im).Error
	return im, err
}

func (h *Handler) GetImage(c *fiber.Ctx) error {
	im, err := h.imageFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "image not found")
	}
	im.ImageURL = h.imageURLFor(im.AssetID)
	return c.JSON(im)
}

func (h *Handler) UpdateImage(c *fiber.Ctx) error {
	im, err := h.imageFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "image not found")
	}
	var in imageInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name != "" {
		im.Name = in.Name
	}
	if in.PushAs != "" {
		im.PushAs = in.PushAs
	}
	if in.FolderID != nil {
		im.FolderID = in.FolderID
	}
	im.Notes = in.Notes
	im.SyncState = "dirty"
	h.db.Save(&im)
	im.ImageURL = h.imageURLFor(im.AssetID)
	return c.JSON(im)
}

func (h *Handler) DeleteImage(c *fiber.Ctx) error {
	im, err := h.imageFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "image not found")
	}
	h.db.Delete(&im)
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) SyncImage(c *fiber.Ctx) error {
	im, err := h.imageFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "image not found")
	}
	conn := h.hub.Get(im.CampaignID)
	if conn == nil {
		return fail(c, fiber.StatusConflict, "foundry not connected")
	}
	img := h.imageURLFor(im.AssetID)
	folder := h.foundryFolderID(im.FolderID)

	var job string
	var doc map[string]any
	if im.PushAs == "journal" {
		job = "create_journal"
		doc = map[string]any{
			"name":  im.Name,
			"pages": []map[string]any{{"name": im.Name, "type": "image", "src": img, "image": map[string]any{"caption": im.Notes}}},
			"flags": map[string]any{"aetherwright": map[string]any{"id": im.ID.String()}},
		}
	} else {
		job = "create_actor"
		doc = map[string]any{
			"name":  im.Name,
			"type":  "npc",
			"img":   img,
			"flags": map[string]any{"aetherwright": map[string]any{"id": im.ID.String()}},
		}
	}
	if folder != "" {
		doc["folder"] = folder
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := conn.Request(ctx, job, doc)
	if err != nil {
		return fail(c, fiber.StatusBadGateway, "foundry: "+err.Error())
	}
	var result struct {
		UUID string `json:"uuid"`
		ID   string `json:"id"`
	}
	_ = json.Unmarshal(resp.Payload, &result)
	im.FoundryUUID = result.UUID
	im.FoundryDocID = result.ID
	im.SyncState = "synced"
	h.db.Save(&im)
	im.ImageURL = img
	return c.JSON(im)
}
