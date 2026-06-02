package api

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aetherwright/backend/internal/crypto"
	"github.com/aetherwright/backend/internal/imagegen"
	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type dalleInput struct {
	Prompt      string     `json:"prompt"`
	Size        string     `json:"size"`
	Provider    string     `json:"provider"`
	SubjectType string     `json:"subjectType"`
	SubjectID   *uuid.UUID `json:"subjectId"`
}

type dalleResult struct {
	ID        uuid.UUID `json:"id"`
	Prompt    string    `json:"prompt"`
	Status    string    `json:"status"`
	Size      string    `json:"size"`
	ImageURL  string    `json:"imageUrl"`
	Mock      bool      `json:"mock"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"createdAt"`
}

func toDalleResult(job models.DalleJob, asset *models.Asset) dalleResult {
	r := dalleResult{
		ID:        job.ID,
		Prompt:    job.Prompt,
		Status:    job.Status,
		Size:      job.Size,
		Error:     job.Error,
		CreatedAt: job.CreatedAt,
	}
	if asset != nil {
		r.ImageURL = asset.StorageURL
		r.Mock = asset.Source == "mock"
	}
	return r
}

func (h *Handler) ListDalle(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var jobs []models.DalleJob
	h.db.Where("campaign_id = ?", cam.ID).Order("created_at desc").Limit(48).Find(&jobs)

	ids := make([]uuid.UUID, 0, len(jobs))
	for _, j := range jobs {
		if j.ResultAssetID != nil {
			ids = append(ids, *j.ResultAssetID)
		}
	}
	assetMap := map[uuid.UUID]models.Asset{}
	if len(ids) > 0 {
		var assets []models.Asset
		h.db.Where("id IN ?", ids).Find(&assets)
		for _, a := range assets {
			assetMap[a.ID] = a
		}
	}

	out := make([]dalleResult, 0, len(jobs))
	for _, j := range jobs {
		var ap *models.Asset
		if j.ResultAssetID != nil {
			if a, ok := assetMap[*j.ResultAssetID]; ok {
				ap = &a
			}
		}
		out = append(out, toDalleResult(j, ap))
	}
	return c.JSON(out)
}

func (h *Handler) GenerateDalle(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in dalleInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	in.Prompt = strings.TrimSpace(in.Prompt)
	if in.Prompt == "" {
		return fail(c, fiber.StatusBadRequest, "prompt required")
	}
	size := in.Size
	if size == "" {
		size = "1024x1024"
	}

	finalPrompt := in.Prompt
	if strings.TrimSpace(cam.StylePrompt) != "" {
		finalPrompt = in.Prompt + " — Stil: " + strings.TrimSpace(cam.StylePrompt)
	}

	job := models.DalleJob{CampaignID: cam.ID, Prompt: finalPrompt, Size: size, Status: "running"}
	h.db.Create(&job)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()
	res, genErr := imagegen.Generate(ctx, h.imageConfigForUser(userID(c), in.Provider), finalPrompt, size)
	if genErr != nil {
		job.Status = "error"
		job.Error = genErr.Error()
		h.db.Save(&job)
		return c.Status(fiber.StatusBadGateway).JSON(toDalleResult(job, nil))
	}
	data, mime, ext, mock := res.Data, res.Mime, res.Ext, res.Mock

	source := "dalle"
	if mock {
		source = "mock"
	}
	asset := models.Asset{CampaignID: cam.ID, Kind: "generated", Source: source, Mime: mime}
	h.db.Create(&asset)

	filename := asset.ID.String() + "." + ext
	if err := os.WriteFile(filepath.Join(h.cfg.UploadDir, filename), data, 0o644); err != nil {
		job.Status = "error"
		job.Error = "could not store image"
		h.db.Save(&job)
		return c.Status(fiber.StatusInternalServerError).JSON(toDalleResult(job, nil))
	}
	asset.StorageURL = strings.TrimRight(h.cfg.PublicBaseURL, "/") + "/uploads/" + filename
	h.db.Save(&asset)

	if in.SubjectID != nil {
		where := "id = ? AND campaign_id = ?"
		switch in.SubjectType {
		case "character":
			h.db.Model(&models.Character{}).Where(where, *in.SubjectID, cam.ID).Update("image_asset_id", asset.ID)
		case "item":
			h.db.Model(&models.Item{}).Where(where, *in.SubjectID, cam.ID).Update("image_asset_id", asset.ID)
		case "scene":
			h.db.Model(&models.Scene{}).Where(where, *in.SubjectID, cam.ID).Update("map_asset_id", asset.ID)
		case "image":
			h.db.Model(&models.ImageEntry{}).Where(where, *in.SubjectID, cam.ID).Update("asset_id", asset.ID)
		}
	}

	job.Status = "done"
	job.ResultAssetID = &asset.ID
	h.db.Save(&job)

	return c.Status(fiber.StatusCreated).JSON(toDalleResult(job, &asset))
}

func (h *Handler) imageConfigForUser(uid uuid.UUID, provider string) imagegen.Config {
	if provider != "" {
		var pc models.ProviderCredential
		if err := h.db.Where("user_id = ? AND provider = ?", uid, provider).First(&pc).Error; err == nil && pc.KeyEnc != "" {
			if key, derr := crypto.Decrypt(h.cfg.EncryptionKey, pc.KeyEnc); derr == nil && key != "" {
				return imagegen.Config{Provider: pc.Provider, Model: pc.Model, APIKey: key}
			}
		}
	}
	if h.cfg.OpenAIKey != "" {
		return imagegen.Config{Provider: "openai", Model: h.cfg.OpenAIModel, APIKey: h.cfg.OpenAIKey}
	}
	return imagegen.Config{}
}
