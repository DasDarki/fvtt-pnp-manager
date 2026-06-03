package api

import (
	"encoding/base64"
	"strings"

	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ExtCampaigns(c *fiber.Ctx) error {
	var rows []models.Campaign
	h.db.Where("owner_id = ?", userID(c)).Order("created_at desc").Find(&rows)
	out := make([]fiber.Map, 0, len(rows))
	for _, r := range rows {
		out = append(out, fiber.Map{"id": r.ID, "name": r.Name, "ruleset": r.Ruleset})
	}
	return c.JSON(out)
}

func (h *Handler) ExtAssetExists(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	ref := c.Query("ref")
	if ref == "" {
		return fail(c, fiber.StatusBadRequest, "ref required")
	}
	var a models.Asset
	if err := h.db.Where("campaign_id = ? AND source_ref = ?", cam.ID, ref).First(&a).Error; err == nil {
		return c.JSON(fiber.Map{"exists": true, "url": a.StorageURL})
	}
	return c.JSON(fiber.Map{"exists": false})
}

type existBatchInput struct {
	Refs []string `json:"refs"`
}

func (h *Handler) ExtAssetsExistBatch(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in existBatchInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if len(in.Refs) == 0 {
		return c.JSON(fiber.Map{"existing": []string{}})
	}
	if len(in.Refs) > 500 {
		in.Refs = in.Refs[:500]
	}
	var rows []models.Asset
	h.db.Where("campaign_id = ? AND source_ref IN ?", cam.ID, in.Refs).Find(&rows)
	existing := make([]string, 0, len(rows))
	for _, a := range rows {
		if a.SourceRef != "" {
			existing = append(existing, a.SourceRef)
		}
	}
	return c.JSON(fiber.Map{"existing": existing})
}

type extUploadInput struct {
	Ref  string `json:"ref"`
	Name string `json:"name"`
	Mime string `json:"mime"`
	Data string `json:"data"`
}

func (h *Handler) ExtUploadAsset(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in extUploadInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Data == "" {
		return fail(c, fiber.StatusBadRequest, "data required")
	}
	if in.Ref != "" {
		var a models.Asset
		if err := h.db.Where("campaign_id = ? AND source_ref = ?", cam.ID, in.Ref).First(&a).Error; err == nil {
			return c.JSON(fiber.Map{"id": a.ID, "url": a.StorageURL, "deduped": true})
		}
	}
	b64 := in.Data
	if strings.HasPrefix(b64, "data:") {
		if i := strings.Index(b64, ","); i >= 0 {
			b64 = b64[i+1:]
		}
	}
	data, derr := base64.StdEncoding.DecodeString(b64)
	if derr != nil {
		return fail(c, fiber.StatusBadRequest, "invalid base64")
	}
	mime, ext := normalizeImageMime(in.Mime)
	asset, serr := h.saveAssetBytes(cam.ID, data, ext, mime, "uploaded", "chatgpt", in.Ref, in.Name)
	if serr != nil {
		return fail(c, fiber.StatusInternalServerError, "could not store file: "+serr.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": asset.ID, "url": asset.StorageURL, "deduped": false})
}
