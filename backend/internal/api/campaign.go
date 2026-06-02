package api

import (
	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type campaignInput struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Ruleset       string `json:"ruleset"`
	ForgeRootPath string `json:"forgeRootPath"`
}

func (h *Handler) campaignFor(c *fiber.Ctx, id string) (models.Campaign, error) {
	var cam models.Campaign
	err := h.db.Where("id = ? AND owner_id = ?", id, userID(c)).First(&cam).Error
	return cam, err
}

func (h *Handler) ListCampaigns(c *fiber.Ctx) error {
	var rows []models.Campaign
	h.db.Where("owner_id = ?", userID(c)).Order("created_at desc").Find(&rows)
	return c.JSON(rows)
}

func (h *Handler) CreateCampaign(c *fiber.Ctx) error {
	var in campaignInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name == "" || in.Ruleset == "" {
		return fail(c, fiber.StatusBadRequest, "name and ruleset required")
	}
	cam := models.Campaign{
		OwnerID:       userID(c),
		Name:          in.Name,
		Slug:          slugify(in.Name),
		Description:   in.Description,
		Ruleset:       in.Ruleset,
		ForgeRootPath: in.ForgeRootPath,
	}
	if err := h.db.Create(&cam).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not create campaign")
	}
	return c.Status(fiber.StatusCreated).JSON(cam)
}

func (h *Handler) GetCampaign(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	return c.JSON(cam)
}

func (h *Handler) UpdateCampaign(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in campaignInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name != "" {
		cam.Name = in.Name
		cam.Slug = slugify(in.Name)
	}
	if in.Description != "" {
		cam.Description = in.Description
	}
	if in.Ruleset != "" {
		cam.Ruleset = in.Ruleset
	}
	if in.ForgeRootPath != "" {
		cam.ForgeRootPath = in.ForgeRootPath
	}
	h.db.Save(&cam)
	return c.JSON(cam)
}

func (h *Handler) DeleteCampaign(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	h.db.Delete(&cam)
	return c.SendStatus(fiber.StatusNoContent)
}

type styleInput struct {
	ArtStyle    string `json:"artStyle"`
	StylePrompt string `json:"stylePrompt"`
}

func (h *Handler) UpdateCampaignStyle(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in styleInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	cam.ArtStyle = in.ArtStyle
	cam.StylePrompt = in.StylePrompt
	h.db.Save(&cam)
	return c.JSON(cam)
}
