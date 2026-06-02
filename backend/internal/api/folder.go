package api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type discoveredFolder struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Parent string `json:"parent"`
}

func (h *Handler) foundryFolderID(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	var f models.Folder
	if err := h.db.First(&f, "id = ?", *id).Error; err != nil {
		return ""
	}
	return f.FoundryFolderID
}

func (h *Handler) ListFolders(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var folders []models.Folder
	h.db.Where("campaign_id = ?", cam.ID).Order("foundry_type asc, sort asc, name asc").Find(&folders)
	return c.JSON(folders)
}

func (h *Handler) DiscoverFolders(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	conn := h.hub.Get(cam.ID)
	if conn == nil {
		return fail(c, fiber.StatusConflict, "foundry not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := conn.Request(ctx, "discover_folders", nil)
	if err != nil {
		return fail(c, fiber.StatusBadGateway, "foundry: "+err.Error())
	}
	var payload struct {
		Folders []discoveredFolder `json:"folders"`
	}
	_ = json.Unmarshal(resp.Payload, &payload)

	idMap := map[string]uuid.UUID{}
	for _, f := range payload.Folders {
		if f.ID == "" {
			continue
		}
		var folder models.Folder
		h.db.Where("campaign_id = ? AND foundry_folder_id = ?", cam.ID, f.ID).First(&folder)
		folder.CampaignID = cam.ID
		folder.FoundryFolderID = f.ID
		folder.Name = f.Name
		folder.FoundryType = f.Type
		folder.Origin = "foundry"
		if folder.ID == uuid.Nil {
			h.db.Create(&folder)
		} else {
			h.db.Save(&folder)
		}
		idMap[f.ID] = folder.ID
	}

	for _, f := range payload.Folders {
		if f.Parent == "" {
			continue
		}
		if pid, ok := idMap[f.Parent]; ok {
			h.db.Model(&models.Folder{}).
				Where("campaign_id = ? AND foundry_folder_id = ?", cam.ID, f.ID).
				Update("parent_id", pid)
		}
	}

	var folders []models.Folder
	h.db.Where("campaign_id = ?", cam.ID).Order("foundry_type asc, name asc").Find(&folders)
	return c.JSON(fiber.Map{"discovered": len(payload.Folders), "folders": folders})
}
