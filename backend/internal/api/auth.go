package api

import (
	"strings"
	"time"

	"github.com/aetherwright/backend/internal/models"
	"github.com/aetherwright/backend/internal/token"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type authInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Device   string `json:"device"`
}

type refreshInput struct {
	RefreshToken string `json:"refreshToken"`
	Device       string `json:"device"`
}

type tokenPair struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	User         models.User `json:"user"`
}

func (h *Handler) issueTokens(u models.User, device string) (tokenPair, error) {
	access, err := token.NewAccessToken(h.cfg.JWTSecret, u.ID, h.cfg.AccessTTL)
	if err != nil {
		return tokenPair{}, err
	}
	raw, hash := token.NewRefreshToken()
	rt := models.RefreshToken{
		UserID:    u.ID,
		TokenHash: hash,
		Device:    device,
		ExpiresAt: time.Now().Add(h.cfg.RefreshTTL),
	}
	if err := h.db.Create(&rt).Error; err != nil {
		return tokenPair{}, err
	}
	return tokenPair{AccessToken: access, RefreshToken: raw, User: u}, nil
}

func (h *Handler) RegisterUser(c *fiber.Ctx) error {
	var in authInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	if in.Email == "" || len(in.Password) < 8 {
		return fail(c, fiber.StatusBadRequest, "email and password (min 8 chars) required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not hash password")
	}
	u := models.User{Email: strings.ToLower(strings.TrimSpace(in.Email)), PasswordHash: string(hash), Name: in.Name}
	if err := h.db.Create(&u).Error; err != nil {
		return fail(c, fiber.StatusConflict, "email already registered")
	}
	pair, err := h.issueTokens(u, in.Device)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not issue tokens")
	}
	return c.Status(fiber.StatusCreated).JSON(pair)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var in authInput
	if err := c.BodyParser(&in); err != nil {
		return fail(c, fiber.StatusBadRequest, "invalid body")
	}
	var u models.User
	if err := h.db.Where("email = ?", strings.ToLower(strings.TrimSpace(in.Email))).First(&u).Error; err != nil {
		return fail(c, fiber.StatusUnauthorized, "invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		return fail(c, fiber.StatusUnauthorized, "invalid credentials")
	}
	pair, err := h.issueTokens(u, in.Device)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not issue tokens")
	}
	return c.JSON(pair)
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	var in refreshInput
	if err := c.BodyParser(&in); err != nil || in.RefreshToken == "" {
		return fail(c, fiber.StatusBadRequest, "refreshToken required")
	}
	var rt models.RefreshToken
	if err := h.db.Where("token_hash = ?", token.HashToken(in.RefreshToken)).First(&rt).Error; err != nil {
		return fail(c, fiber.StatusUnauthorized, "invalid refresh token")
	}
	if rt.RevokedAt != nil || time.Now().After(rt.ExpiresAt) {
		return fail(c, fiber.StatusUnauthorized, "refresh token expired")
	}
	now := time.Now()
	rt.RevokedAt = &now
	h.db.Save(&rt)

	var u models.User
	if err := h.db.First(&u, "id = ?", rt.UserID).Error; err != nil {
		return fail(c, fiber.StatusUnauthorized, "user not found")
	}
	device := in.Device
	if device == "" {
		device = rt.Device
	}
	pair, err := h.issueTokens(u, device)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "could not issue tokens")
	}
	return c.JSON(pair)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	var in refreshInput
	if err := c.BodyParser(&in); err != nil || in.RefreshToken == "" {
		return fail(c, fiber.StatusBadRequest, "refreshToken required")
	}
	h.db.Model(&models.RefreshToken{}).
		Where("token_hash = ?", token.HashToken(in.RefreshToken)).
		Update("revoked_at", time.Now())
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) Me(c *fiber.Ctx) error {
	var u models.User
	if err := h.db.First(&u, "id = ?", userID(c)).Error; err != nil {
		return fail(c, fiber.StatusNotFound, "user not found")
	}
	return c.JSON(u)
}
