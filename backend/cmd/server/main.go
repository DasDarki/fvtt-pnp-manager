package main

import (
	"log"
	"os"

	"github.com/aetherwright/backend/internal/api"
	"github.com/aetherwright/backend/internal/config"
	"github.com/aetherwright/backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connect: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migrate: %v", err)
	}

	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatalf("upload dir: %v", err)
	}

	app := fiber.New(fiber.Config{AppName: "Aetherwright API", BodyLimit: 16 * 1024 * 1024})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	// Uploaded images are public assets; allow any origin to read them so the
	// FoundryVTT / The Forge instance can fetch the bytes (to re-host them) and
	// load them as WebGL token textures without cross-origin tainting.
	app.Use("/uploads", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Cross-Origin-Resource-Policy", "cross-origin")
		return c.Next()
	})
	app.Static("/uploads", cfg.UploadDir, fiber.Static{ByteRange: true})

	api.New(db, cfg).Routes(app)

	log.Printf("Aetherwright API listening on :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
