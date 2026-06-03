package api

import (
	"strings"
	"time"

	"github.com/aetherwright/backend/internal/models"
	"github.com/aetherwright/backend/internal/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return fail(c, fiber.StatusUnauthorized, "missing bearer token")
		}
		uid, err := token.ParseAccessToken(h.cfg.JWTSecret, strings.TrimPrefix(auth, "Bearer "))
		if err != nil || uid == uuid.Nil {
			return fail(c, fiber.StatusUnauthorized, "invalid token")
		}
		c.Locals("userID", uid)
		return c.Next()
	}
}

func userID(c *fiber.Ctx) uuid.UUID {
	return c.Locals("userID").(uuid.UUID)
}

func (h *Handler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Get("X-API-Key")
		if raw == "" {
			if auth := c.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
				raw = strings.TrimPrefix(auth, "Bearer ")
			}
		}
		if !strings.HasPrefix(raw, "awk_") {
			return fail(c, fiber.StatusUnauthorized, "missing api key")
		}
		var key models.ApiKey
		if err := h.db.Where("key_hash = ?", token.HashToken(raw)).First(&key).Error; err != nil {
			return fail(c, fiber.StatusUnauthorized, "invalid api key")
		}
		c.Locals("userID", key.UserID)
		c.Locals("apiKeyScope", key.Scope)
		now := time.Now()
		h.db.Model(&models.ApiKey{}).Where("id = ?", key.ID).Update("last_used_at", now)
		return c.Next()
	}
}
