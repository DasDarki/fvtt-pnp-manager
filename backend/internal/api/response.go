package api

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func fail(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(fiber.Map{"error": msg})
}

var slugStrip = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	out := slugStrip.ReplaceAllString(strings.ToLower(s), "-")
	return strings.Trim(out, "-")
}
