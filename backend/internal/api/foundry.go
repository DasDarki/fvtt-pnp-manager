package api

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/aetherwright/backend/internal/models"
	"github.com/aetherwright/backend/internal/token"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) wsURL() string {
	return strings.Replace(strings.TrimRight(h.cfg.PublicBaseURL, "/"), "http", "ws", 1) + "/ws/foundry"
}

func (h *Handler) FoundryWS(ws *websocket.Conn) {
	tok := ws.Query("token")
	if tok == "" {
		_ = ws.Close()
		return
	}
	var fc models.FoundryConnection
	if err := h.db.Where("pairing_token = ?", token.HashToken(tok)).First(&fc).Error; err != nil {
		_ = ws.Close()
		return
	}
	now := time.Now()
	h.db.Model(&models.FoundryConnection{}).Where("id = ?", fc.ID).
		Updates(map[string]any{"status": "connected", "last_seen_at": now})

	h.hub.Serve(ws, fc.CampaignID, func() {
		h.db.Model(&models.FoundryConnection{}).Where("id = ?", fc.ID).Update("status", "disconnected")
	})
}

func (h *Handler) PairFoundry(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	raw, hash := token.NewRefreshToken()

	var fc models.FoundryConnection
	if err := h.db.Where("campaign_id = ?", cam.ID).First(&fc).Error; err != nil {
		fc = models.FoundryConnection{CampaignID: cam.ID, Label: "Foundry", PairingToken: hash, Status: "disconnected"}
		h.db.Create(&fc)
	} else {
		fc.PairingToken = hash
		fc.Status = "disconnected"
		h.db.Save(&fc)
	}
	return c.JSON(fiber.Map{"token": raw, "wsUrl": h.wsURL()})
}

func (h *Handler) FoundryStatus(c *fiber.Ctx) error {
	cam, err := h.campaignFor(c, c.Params("id"))
	if err != nil {
		return fail(c, fiber.StatusNotFound, "campaign not found")
	}
	var fc models.FoundryConnection
	paired := h.db.Where("campaign_id = ?", cam.ID).First(&fc).Error == nil
	conn := h.hub.Get(cam.ID)
	resp := fiber.Map{
		"paired":     paired,
		"connected":  conn != nil,
		"wsUrl":      h.wsURL(),
		"lastSeenAt": fc.LastSeenAt,
	}
	if conn != nil {
		resp["world"] = conn.World
		resp["version"] = conn.Version
	}
	return c.JSON(resp)
}

func (h *Handler) SyncCharacter(c *fiber.Ctx) error {
	ch, err := h.characterFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "character not found")
	}
	conn := h.hub.Get(ch.CampaignID)
	if conn == nil {
		return fail(c, fiber.StatusConflict, "foundry not connected")
	}

	doc := h.toFoundryActor(ch)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := conn.Request(ctx, "create_actor", doc)
	if err != nil {
		return fail(c, fiber.StatusBadGateway, "foundry: "+err.Error())
	}
	var result struct {
		UUID string `json:"uuid"`
		ID   string `json:"id"`
		Img  string `json:"img"`
	}
	_ = json.Unmarshal(resp.Payload, &result)

	ch.FoundryUUID = result.UUID
	ch.FoundryDocID = result.ID
	ch.SyncState = "synced"
	h.db.Save(&ch)
	ch.ImageURL = h.imageURLFor(ch.ImageAssetID)
	return c.JSON(ch)
}

func (h *Handler) toFoundryActor(ch models.Character) map[string]any {
	var sd map[string]any
	_ = json.Unmarshal(ch.SystemData, &sd)
	if sd == nil {
		sd = map[string]any{}
	}
	ab := subMap(sd, "abilities")
	saves := subMap(sd, "saves")
	sk := subMap(sd, "skills")

	abilities := map[string]any{}
	for _, k := range []string{"str", "dex", "con", "int", "wis", "cha"} {
		abilities[k] = map[string]any{"value": numOf(ab, k), "proficient": numOf(saves, k)}
	}
	skills := map[string]any{}
	for _, k := range []string{"acr", "ani", "arc", "ath", "dec", "his", "ins", "itm", "inv", "med", "nat", "prc", "prf", "per", "rel", "slt", "ste", "sur"} {
		skills[k] = map[string]any{"value": numOf(sk, k)}
	}

	actorType := "npc"
	if ch.CharacterType == "pc" || ch.CharacterType == "ally" {
		actorType = "character"
	}

	doc := map[string]any{
		"name": ch.Name,
		"type": actorType,
		"img":  h.imageURLFor(ch.ImageAssetID),
		"system": map[string]any{
			"abilities": abilities,
			"skills":    skills,
			"attributes": map[string]any{
				"ac":           map[string]any{"flat": numOf(sd, "ac"), "calc": "flat"},
				"hp":           map[string]any{"value": numOf(sd, "hp"), "max": numOf(sd, "hpMax"), "temp": numOf(sd, "hpTemp")},
				"init":         map[string]any{"bonus": numOf(sd, "initiative")},
				"movement":     map[string]any{"walk": numOf(sd, "speed")},
				"spellcasting": strOf(sd, "spellcasting"),
			},
			"details": map[string]any{
				"alignment":  strOf(sd, "alignment"),
				"race":       strOf(sd, "race"),
				"background": strOf(sd, "background"),
				"level":      numOf(sd, "level"),
				"biography":  map[string]any{"value": strOf(sd, "notes")},
			},
			"currency": map[string]any{
				"pp": numOf(sd, "pp"), "gp": numOf(sd, "gp"), "ep": numOf(sd, "ep"), "sp": numOf(sd, "sp"), "cp": numOf(sd, "cp"),
			},
		},
		"flags": map[string]any{
			"aetherwright": map[string]any{"id": ch.ID.String(), "class": strOf(sd, "class")},
		},
	}
	if fid := h.foundryFolderID(ch.FolderID); fid != "" {
		doc["folder"] = fid
	}
	return doc
}

func (h *Handler) SyncItem(c *fiber.Ctx) error {
	it, err := h.itemFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "item not found")
	}
	conn := h.hub.Get(it.CampaignID)
	if conn == nil {
		return fail(c, fiber.StatusConflict, "foundry not connected")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := conn.Request(ctx, "create_item", h.toFoundryItem(it))
	if err != nil {
		return fail(c, fiber.StatusBadGateway, "foundry: "+err.Error())
	}
	var result struct {
		UUID string `json:"uuid"`
		ID   string `json:"id"`
	}
	_ = json.Unmarshal(resp.Payload, &result)
	it.FoundryUUID = result.UUID
	it.FoundryDocID = result.ID
	it.SyncState = "synced"
	h.db.Save(&it)
	it.ImageURL = h.imageURLFor(it.ImageAssetID)
	return c.JSON(it)
}

func (h *Handler) toFoundryItem(it models.Item) map[string]any {
	var sd map[string]any
	_ = json.Unmarshal(it.SystemData, &sd)
	if sd == nil {
		sd = map[string]any{}
	}
	kind := strOf(sd, "kind")
	if kind == "" {
		kind = "loot"
	}
	attunement := ""
	if it.Attuned {
		attunement = "required"
	}
	doc := map[string]any{
		"name": it.Name,
		"type": kind,
		"img":  h.imageURLFor(it.ImageAssetID),
		"system": map[string]any{
			"description": map[string]any{"value": it.Summary},
			"rarity":      it.Rarity,
			"quantity":    1,
			"weight":      map[string]any{"value": numOf(sd, "weight"), "units": "lb"},
			"price":       map[string]any{"value": numOf(sd, "value"), "denomination": "gp"},
			"attunement":  attunement,
			"attuned":     it.Attuned,
		},
		"flags": map[string]any{"aetherwright": map[string]any{"id": it.ID.String()}},
	}
	if fid := h.foundryFolderID(it.FolderID); fid != "" {
		doc["folder"] = fid
	}
	return doc
}

func (h *Handler) SyncScene(c *fiber.Ctx) error {
	s, err := h.sceneFor(c)
	if err != nil {
		return fail(c, fiber.StatusNotFound, "scene not found")
	}
	conn := h.hub.Get(s.CampaignID)
	if conn == nil {
		return fail(c, fiber.StatusConflict, "foundry not connected")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := conn.Request(ctx, "create_scene", h.toFoundryScene(s))
	if err != nil {
		return fail(c, fiber.StatusBadGateway, "foundry: "+err.Error())
	}
	var result struct {
		UUID string `json:"uuid"`
		ID   string `json:"id"`
	}
	_ = json.Unmarshal(resp.Payload, &result)
	s.FoundryUUID = result.UUID
	s.FoundryDocID = result.ID
	s.SyncState = "synced"
	h.db.Save(&s)
	s.ImageURL = h.imageURLFor(s.MapAssetID)
	return c.JSON(s)
}

func (h *Handler) toFoundryScene(s models.Scene) map[string]any {
	var sd map[string]any
	_ = json.Unmarshal(s.SystemData, &sd)
	if sd == nil {
		sd = map[string]any{}
	}
	doc := map[string]any{
		"name":       s.Name,
		"navigation": true,
		"navName":    strOf(sd, "act"),
		"flags": map[string]any{
			"aetherwright": map[string]any{"id": s.ID.String(), "summary": s.Summary, "status": strOf(sd, "status")},
		},
	}
	if img := h.imageURLFor(s.MapAssetID); img != "" {
		doc["background"] = map[string]any{"src": img}
	}
	if fid := h.foundryFolderID(s.FolderID); fid != "" {
		doc["folder"] = fid
	}
	return doc
}

func numOf(m map[string]any, k string) float64 {
	if m == nil {
		return 0
	}
	if v, ok := m[k].(float64); ok {
		return v
	}
	return 0
}

func strOf(m map[string]any, k string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[k].(string); ok {
		return v
	}
	return ""
}

func subMap(m map[string]any, k string) map[string]any {
	if m == nil {
		return nil
	}
	if v, ok := m[k].(map[string]any); ok {
		return v
	}
	return nil
}
