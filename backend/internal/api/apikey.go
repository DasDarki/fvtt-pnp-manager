package api

import (
	"strings"
	"time"

	"github.com/aetherwright/backend/internal/models"
	"github.com/aetherwright/backend/internal/token"
	"github.com/gofiber/fiber/v2"
)

type apiKeyOut struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	Scope      string     `json:"scope"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time  `json:"createdAt"`
}

func (h *Handler) ListApiKeys(c *fiber.Ctx) error {
	var keys []models.ApiKey
	h.db.Where("user_id = ?", userID(c)).Order("created_at desc").Find(&keys)
	out := make([]apiKeyOut, 0, len(keys))
	for _, k := range keys {
		out = append(out, apiKeyOut{ID: k.ID.String(), Name: k.Name, Prefix: k.Prefix, Scope: k.Scope, LastUsedAt: k.LastUsedAt, CreatedAt: k.CreatedAt})
	}
	return c.JSON(out)
}

type apiKeyInput struct {
	Name string `json:"name"`
}

func (h *Handler) CreateApiKey(c *fiber.Ctx) error {
	var in apiKeyInput
	_ = c.BodyParser(&in)
	name := strings.TrimSpace(in.Name)
	if name == "" {
		name = "ChatGPT Extension"
	}
	raw, _ := token.NewRefreshToken()
	full := "awk_" + raw
	key := models.ApiKey{
		UserID:  userID(c),
		Name:    name,
		Prefix:  full[:12],
		KeyHash: token.HashToken(full),
		Scope:   "images",
	}
	h.db.Create(&key)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":     key.ID.String(),
		"name":   key.Name,
		"prefix": key.Prefix,
		"scope":  key.Scope,
		"key":    full,
	})
}

func (h *Handler) DeleteApiKey(c *fiber.Ctx) error {
	h.db.Where("id = ? AND user_id = ?", c.Params("id"), userID(c)).Delete(&models.ApiKey{})
	return c.SendStatus(fiber.StatusNoContent)
}
