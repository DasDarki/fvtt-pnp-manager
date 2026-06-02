package api

import (
	"strings"

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
