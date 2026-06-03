package api

import (
	"github.com/aetherwright/backend/internal/config"
	"github.com/aetherwright/backend/internal/foundry"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	cfg config.Config
	hub *foundry.Hub
}

func New(db *gorm.DB, cfg config.Config) *Handler {
	return &Handler{db: db, cfg: cfg, hub: foundry.NewHub()}
}

func (h *Handler) Routes(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/foundry", websocket.New(h.FoundryWS))

	app.Get("/health", h.health)

	// The API is mounted under both /api/v1 and /v1. A reverse proxy (Coolify /
	// Traefik) that routes the backend on a /api path prefix and strips it will
	// forward requests as /v1/...; direct access keeps working under /api/v1/...
	h.registerAPI(app.Group("/api/v1"))
	h.registerAPI(app.Group("/v1"))
}

func (h *Handler) health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *Handler) registerAPI(api fiber.Router) {
	api.Get("/health", h.health)

	auth := api.Group("/auth")
	auth.Post("/register", h.RegisterUser)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)
	auth.Get("/me", h.Protected(), h.Me)

	// API-key authenticated (browser extension). Registered before the JWT
	// catch-all group below so the Protected() Use does not apply to it.
	ext := api.Group("/ext", h.ApiKeyAuth())
	ext.Get("/campaigns", h.ExtCampaigns)
	ext.Get("/campaigns/:campaignId/assets/exists", h.ExtAssetExists)
	ext.Post("/campaigns/:campaignId/assets/exists", h.ExtAssetsExistBatch)
	ext.Post("/campaigns/:campaignId/assets", h.ExtUploadAsset)

	p := api.Group("", h.Protected())

	p.Get("/providers", h.ListProviders)
	p.Put("/providers/:provider", h.SetProvider)
	p.Delete("/providers/:provider", h.DeleteProvider)

	p.Get("/apikeys", h.ListApiKeys)
	p.Post("/apikeys", h.CreateApiKey)
	p.Delete("/apikeys/:id", h.DeleteApiKey)

	p.Get("/campaigns", h.ListCampaigns)
	p.Post("/campaigns", h.CreateCampaign)
	p.Get("/campaigns/:id/stats", h.CampaignStats)
	p.Get("/campaigns/:id", h.GetCampaign)
	p.Patch("/campaigns/:id", h.UpdateCampaign)
	p.Patch("/campaigns/:id/style", h.UpdateCampaignStyle)
	p.Post("/campaigns/:id/foundry/pair", h.PairFoundry)
	p.Get("/campaigns/:id/foundry", h.FoundryStatus)
	p.Post("/campaigns/:id/foundry/discover", h.DiscoverFolders)
	p.Get("/campaigns/:id/folders", h.ListFolders)
	p.Post("/campaigns/:id/folders", h.CreateFolder)
	p.Delete("/campaigns/:id/folders/:folderId", h.DeleteFolder)
	p.Delete("/campaigns/:id", h.DeleteCampaign)

	p.Get("/campaigns/:campaignId/characters", h.ListCharacters)
	p.Post("/campaigns/:campaignId/characters", h.CreateCharacter)
	p.Get("/campaigns/:campaignId/characters/:id", h.GetCharacter)
	p.Patch("/campaigns/:campaignId/characters/:id", h.UpdateCharacter)
	p.Delete("/campaigns/:campaignId/characters/:id", h.DeleteCharacter)
	p.Post("/campaigns/:campaignId/characters/:id/sync", h.SyncCharacter)

	p.Get("/campaigns/:campaignId/items", h.ListItems)
	p.Post("/campaigns/:campaignId/items", h.CreateItem)
	p.Get("/campaigns/:campaignId/items/:id", h.GetItem)
	p.Patch("/campaigns/:campaignId/items/:id", h.UpdateItem)
	p.Delete("/campaigns/:campaignId/items/:id", h.DeleteItem)
	p.Post("/campaigns/:campaignId/items/:id/sync", h.SyncItem)

	p.Get("/campaigns/:campaignId/scenes", h.ListScenes)
	p.Post("/campaigns/:campaignId/scenes", h.CreateScene)
	p.Get("/campaigns/:campaignId/scenes/:id", h.GetScene)
	p.Patch("/campaigns/:campaignId/scenes/:id", h.UpdateScene)
	p.Delete("/campaigns/:campaignId/scenes/:id", h.DeleteScene)
	p.Post("/campaigns/:campaignId/scenes/:id/sync", h.SyncScene)

	p.Get("/campaigns/:campaignId/images", h.ListImages)
	p.Post("/campaigns/:campaignId/images", h.CreateImage)
	p.Get("/campaigns/:campaignId/images/:id", h.GetImage)
	p.Patch("/campaigns/:campaignId/images/:id", h.UpdateImage)
	p.Delete("/campaigns/:campaignId/images/:id", h.DeleteImage)
	p.Post("/campaigns/:campaignId/images/:id/sync", h.SyncImage)

	p.Get("/campaigns/:campaignId/memories", h.ListMemories)
	p.Post("/campaigns/:campaignId/memories", h.CreateMemory)
	p.Patch("/campaigns/:campaignId/memories/:id", h.UpdateMemory)
	p.Delete("/campaigns/:campaignId/memories/:id", h.DeleteMemory)
	p.Get("/campaigns/:campaignId/alerts", h.ListAlerts)

	p.Get("/campaigns/:campaignId/links", h.ListLinks)
	p.Post("/campaigns/:campaignId/links", h.CreateLink)
	p.Delete("/campaigns/:campaignId/links/:id", h.DeleteLink)

	p.Get("/campaigns/:campaignId/tags", h.ListTags)
	p.Post("/campaigns/:campaignId/tags", h.CreateTag)
	p.Delete("/campaigns/:campaignId/tags/:id", h.DeleteTag)
	p.Get("/campaigns/:campaignId/entity-tags", h.ListEntityTags)
	p.Post("/campaigns/:campaignId/entity-tags", h.CreateEntityTag)
	p.Delete("/campaigns/:campaignId/entity-tags/:id", h.DeleteEntityTag)

	p.Get("/campaigns/:campaignId/dalle", h.ListDalle)
	p.Post("/campaigns/:campaignId/dalle/generate", h.GenerateDalle)

	p.Get("/campaigns/:campaignId/assets", h.ListAssets)
	p.Post("/campaigns/:campaignId/assets", h.UploadAsset)
	p.Post("/campaigns/:campaignId/assets/attach", h.AttachAsset)
	p.Patch("/campaigns/:campaignId/assets/:id", h.RenameAsset)
	p.Delete("/campaigns/:campaignId/assets/:id", h.DeleteAsset)
}
