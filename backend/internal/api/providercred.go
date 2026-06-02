package api

import (
	"strings"

	"github.com/aetherwright/backend/internal/crypto"
	"github.com/aetherwright/backend/internal/imagegen"
	"github.com/aetherwright/backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type providerCredInput struct {
	Model  string `json:"model"`
	APIKey string `json:"apiKey"`
}

type providerCredOut struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	HasKey   bool   `json:"hasKey"`
}

func providerInCatalog(p string) bool {
	for _, x := range imagegen.Providers {
		if x == p {
			return true
		}
	}
	return false
}

func (h *Handler) ListProviders(c *fiber.Ctx) error {
	var rows []models.ProviderCredential
	h.db.Where("user_id = ?", userID(c)).Order("provider asc").Find(&rows)
	out := make([]providerCredOut, 0, len(rows))
	for _, r := range rows {
		out = append(out, providerCredOut{Provider: r.Provider, Model: r.Model, HasKey: r.KeyEnc != ""})
	}
	return c.JSON(fiber.Map{"providers": out, "catalog": imagegen.Providers})
}

func (h *Handler) SetProvider(c *fiber.Ctx) error {
	prov := c.Params("provider")
	if !providerInCatalog(prov) {
		return fail(c, fiber.StatusBadRequest, "unknown provider")
	}
	var in providerCredInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	uid := userID(c)
	var pc models.ProviderCredential
	if err := h.db.Where("user_id = ? AND provider = ?", uid, prov).First(&pc).Error; err != nil {
		pc = models.ProviderCredential{UserID: uid, Provider: prov}
	}
	pc.Model = in.Model
	if strings.TrimSpace(in.APIKey) != "" {
		enc, eerr := crypto.Encrypt(h.cfg.EncryptionKey, strings.TrimSpace(in.APIKey))
		if eerr != nil {
			return fail(c, fiber.StatusInternalServerError, "could not store key")
		}
		pc.KeyEnc = enc
	}
	if pc.ID == uuid.Nil {
		h.db.Create(&pc)
	} else {
		h.db.Save(&pc)
	}
	return c.JSON(providerCredOut{Provider: pc.Provider, Model: pc.Model, HasKey: pc.KeyEnc != ""})
}

func (h *Handler) DeleteProvider(c *fiber.Ctx) error {
	h.db.Where("user_id = ? AND provider = ?", userID(c), c.Params("provider")).Delete(&models.ProviderCredential{})
	return c.SendStatus(fiber.StatusNoContent)
}
