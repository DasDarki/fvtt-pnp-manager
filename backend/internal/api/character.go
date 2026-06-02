package api

import (
	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type characterInput struct {
	Name          string         `json:"name"`
	Summary       string         `json:"summary"`
	CharacterType string         `json:"characterType"`
	Status        string         `json:"status"`
	SystemData    datatypes.JSON `json:"systemData"`
	FolderID      *uuid.UUID     `json:"folderId"`
}

func (h *Handler) imageURLFor(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	var a models.Asset
	if err := h.db.First(&a, "id = ?", *id).Error; err != nil {
		return ""
	}
	return a.StorageURL
}

func (h *Handler) ListCharacters(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var rows []models.Character
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

func (h *Handler) CreateCharacter(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var in characterInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name == "" {
		return fail(c, fiber.StatusBadRequest, "name required")
	}
	ch := models.Character{
		CampaignID:    cam.ID,
		Name:          in.Name,
		Summary:       in.Summary,
		CharacterType: orDefault(in.CharacterType, "npc"),
		Status:        orDefault(in.Status, "alive"),
		SystemData:    in.SystemData,
		FolderID:      in.FolderID,
	}
	if err := h.db.Create(&ch).Error; err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not create character")
	}
	return c.Status(fiber.StatusCreated).JSON(ch)
}

func (h *Handler) characterFor(c *fiber.Ctx) (models.Character, error) {
	cam, err := h.campaignFor(c, c.Params("campaignId"))
	if err != nil {
		return models.Character{}, err
	}
	var ch models.Character
	err = h.db.Where("id = ? AND campaign_id = ?", c.Params("id"), cam.ID).First(&ch).Error
	return ch, err
}

func (h *Handler) GetCharacter(c *fiber.Ctx) error {
	ch, err := h.characterFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "character not found")
	}
	ch.ImageURL = h.imageURLFor(ch.ImageAssetID)
	return c.JSON(ch)
}

func (h *Handler) UpdateCharacter(c *fiber.Ctx) error {
	ch, err := h.characterFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "character not found")
	}
	var in characterInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Name != "" {
		ch.Name = in.Name
	}
	if in.Summary != "" {
		ch.Summary = in.Summary
	}
	if in.CharacterType != "" {
		ch.CharacterType = in.CharacterType
	}
	if in.Status != "" {
		ch.Status = in.Status
	}
	if in.SystemData != nil {
		ch.SystemData = in.SystemData
	}
	if in.FolderID != nil {
		ch.FolderID = in.FolderID
	}
	ch.SyncState = "dirty"
	h.db.Save(&ch)
	return c.JSON(ch)
}

func (h *Handler) DeleteCharacter(c *fiber.Ctx) error {
	ch, err := h.characterFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "character not found")
	}
	h.db.Delete(&ch)
	return c.SendStatus(fiber.StatusNoContent)
}

func orDefault(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
