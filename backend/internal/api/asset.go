package api

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type assetOut struct {
	ID        uuid.UUID `json:"id"`
	URL       string    `json:"url"`
	Mime      string    `json:"mime"`
	Kind      string    `json:"kind"`
	Source    string    `json:"source"`
	Prompt    string    `json:"prompt"`
	CreatedAt time.Time `json:"createdAt"`
}

var uploadExt = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".webp": "image/webp",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
}

func (h *Handler) toAssetOut(a models.Asset, prompt string) assetOut {
	return assetOut{ID: a.ID, URL: a.StorageURL, Mime: a.Mime, Kind: a.Kind, Source: a.Source, Prompt: prompt, CreatedAt: a.CreatedAt}
}

func normalizeImageMime(mime string) (string, string) {
	switch {
	case strings.Contains(mime, "png"):
		return "image/png", "png"
	case strings.Contains(mime, "jpeg"), strings.Contains(mime, "jpg"):
		return "image/jpeg", "jpg"
	case strings.Contains(mime, "webp"):
		return "image/webp", "webp"
	case strings.Contains(mime, "gif"):
		return "image/gif", "gif"
	case strings.Contains(mime, "svg"):
		return "image/svg+xml", "svg"
	default:
		return "image/png", "png"
	}
}

func (h *Handler) saveAssetBytes(camID uuid.UUID, data []byte, ext, mime, kind, source, ref string) (models.Asset, error) {
	asset := models.Asset{CampaignID: camID, Kind: kind, Source: source, Mime: mime, SourceRef: ref}
	h.db.Create(&asset)
	filename := asset.ID.String() + "." + ext
	if err := os.WriteFile(filepath.Join(h.cfg.UploadDir, filename), data, 0o644); err != nil {
		h.db.Delete(&asset)
		return models.Asset{}, err
	}
	asset.StorageURL = strings.TrimRight(h.cfg.PublicBaseURL, "/") + "/uploads/" + filename
	h.db.Save(&asset)
	return asset, nil
}

func (h *Handler) ListAssets(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var assets []models.Asset
	h.db.Where("campaign_id = ? AND storage_url <> ''", cam.ID).Order("created_at desc").Limit(500).Find(&assets)

	prompts := map[uuid.UUID]string{}
	var jobs []models.DalleJob
	h.db.Where("campaign_id = ? AND result_asset_id IS NOT NULL", cam.ID).Find(&jobs)
	for _, j := range jobs {
		if j.ResultAssetID != nil {
			prompts[*j.ResultAssetID] = j.Prompt
		}
	}

	out := make([]assetOut, 0, len(assets))
	for _, a := range assets {
		out = append(out, h.toAssetOut(a, prompts[a.ID]))
	}
	return c.JSON(out)
}

func (h *Handler) UploadAsset(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "no file")
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	mime, ok := uploadExt[ext]
	if !ok {
		return fail(c, fiber.StatusBadRequest, "unsupported file type")
	}

	asset := models.Asset{CampaignID: cam.ID, Kind: "uploaded", Source: "upload", Mime: mime}
	h.db.Create(&asset)

	filename := asset.ID.String() + ext
	if err := c.SaveFile(fh, filepath.Join(h.cfg.UploadDir, filename)); err != nil {
		h.db.Delete(&asset)
		return fail(c, fiber.StatusInternalServerError, "could not store file: "+err.Error())
	}
	asset.StorageURL = strings.TrimRight(h.cfg.PublicBaseURL, "/") + "/uploads/" + filename
	h.db.Save(&asset)
	return c.Status(fiber.StatusCreated).JSON(h.toAssetOut(asset, ""))
}

func (h *Handler) DeleteAsset(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var a models.Asset
	if err := h.db.Where("id = ? AND campaign_id = ?", c.Params("id"), cam.ID).First(&a).Error; err != nil {
		return fail(c, fiber.StatusNotFound, "asset not found")
	}
	if i := strings.LastIndex(a.StorageURL, "/uploads/"); i >= 0 {
		_ = os.Remove(filepath.Join(h.cfg.UploadDir, a.StorageURL[i+len("/uploads/"):]))
	}
	h.db.Delete(&a)
	return c.SendStatus(fiber.StatusNoContent)
}

type attachInput struct {
	SubjectType string    `json:"subjectType"`
	SubjectID   uuid.UUID `json:"subjectId"`
	AssetID     uuid.UUID `json:"assetId"`
}

func (h *Handler) AttachAsset(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in attachInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	var a models.Asset
	if err := h.db.Where("id = ? AND campaign_id = ?", in.AssetID, cam.ID).First(&a).Error; err != nil {
		return fail(c, fiber.StatusNotFound, "asset not found")
	}
	where := "id = ? AND campaign_id = ?"
	switch in.SubjectType {
	case "character":
		h.db.Model(&models.Character{}).Where(where, in.SubjectID, cam.ID).Updates(map[string]any{"image_asset_id": a.ID, "sync_state": "dirty"})
	case "item":
		h.db.Model(&models.Item{}).Where(where, in.SubjectID, cam.ID).Updates(map[string]any{"image_asset_id": a.ID, "sync_state": "dirty"})
	case "scene":
		h.db.Model(&models.Scene{}).Where(where, in.SubjectID, cam.ID).Updates(map[string]any{"map_asset_id": a.ID, "sync_state": "dirty"})
	case "image":
		h.db.Model(&models.ImageEntry{}).Where(where, in.SubjectID, cam.ID).Updates(map[string]any{"asset_id": a.ID, "sync_state": "dirty"})
	default:
		return fail(c, fiber.StatusBadRequest, "unknown subject type")
	}
	return c.JSON(fiber.Map{"imageUrl": a.StorageURL})
}
