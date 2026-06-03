package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Base) BeforeCreate(*gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type User struct {
	Base
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Name         string `json:"name"`
}

type RefreshToken struct {
	Base
	UserID    uuid.UUID  `gorm:"type:uuid;index;not null" json:"-"`
	TokenHash string     `gorm:"uniqueIndex;not null" json:"-"`
	Device    string     `json:"device"`
	ExpiresAt time.Time  `json:"expiresAt"`
	RevokedAt *time.Time `json:"-"`
}

type Campaign struct {
	Base
	OwnerID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"-"`
	Name          string         `gorm:"not null" json:"name"`
	Slug          string         `json:"slug"`
	Description   string         `json:"description"`
	Ruleset       string         `gorm:"not null" json:"ruleset"`
	ForgeRootPath string         `json:"forgeRootPath"`
	ArtStyle      string         `json:"artStyle"`
	StylePrompt   string         `json:"stylePrompt"`
	Settings      datatypes.JSON `gorm:"type:jsonb" json:"settings"`
}

type ProviderCredential struct {
	Base
	UserID   uuid.UUID `gorm:"type:uuid;index;not null" json:"-"`
	Provider string    `gorm:"index" json:"provider"`
	Model    string    `json:"model"`
	KeyEnc   string    `json:"-"`
}

type ApiKey struct {
	Base
	UserID     uuid.UUID  `gorm:"type:uuid;index;not null" json:"-"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	KeyHash    string     `gorm:"uniqueIndex" json:"-"`
	Scope      string     `gorm:"default:images" json:"scope"`
	LastUsedAt *time.Time `json:"lastUsedAt"`
}

type Folder struct {
	Base
	CampaignID      uuid.UUID  `gorm:"type:uuid;index;not null" json:"campaignId"`
	ParentID        *uuid.UUID `gorm:"type:uuid;index" json:"parentId"`
	Name            string     `json:"name"`
	Color           string     `json:"color"`
	Sort            int        `json:"sort"`
	FoundryFolderID string     `json:"foundryFolderId"`
	FoundryType     string     `json:"foundryType"`
	Origin          string     `gorm:"default:local" json:"origin"`
}

type Character struct {
	Base
	CampaignID    uuid.UUID      `gorm:"type:uuid;index;not null" json:"campaignId"`
	FolderID      *uuid.UUID     `gorm:"type:uuid;index" json:"folderId"`
	Name          string         `gorm:"not null" json:"name"`
	Summary       string         `json:"summary"`
	CharacterType string         `gorm:"default:npc" json:"characterType"`
	Status        string         `gorm:"default:alive" json:"status"`
	SystemData    datatypes.JSON `gorm:"type:jsonb" json:"systemData"`
	ImageAssetID  *uuid.UUID     `gorm:"type:uuid" json:"imageAssetId"`
	ImageURL      string         `gorm:"-" json:"imageUrl"`
	FoundryUUID   string         `json:"foundryUuid"`
	FoundryDocID  string         `json:"foundryDocId"`
	SyncState     string         `gorm:"default:none" json:"syncState"`
	Sort          int            `json:"sort"`
}

type Item struct {
	Base
	CampaignID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"campaignId"`
	FolderID     *uuid.UUID     `gorm:"type:uuid;index" json:"folderId"`
	Name         string         `gorm:"not null" json:"name"`
	Summary      string         `json:"summary"`
	ItemType     string         `json:"itemType"`
	Rarity       string         `gorm:"default:common" json:"rarity"`
	Attuned      bool           `json:"attuned"`
	SystemData   datatypes.JSON `gorm:"type:jsonb" json:"systemData"`
	ImageAssetID *uuid.UUID     `gorm:"type:uuid" json:"imageAssetId"`
	ImageURL     string         `gorm:"-" json:"imageUrl"`
	FoundryUUID  string         `json:"foundryUuid"`
	FoundryDocID string         `json:"foundryDocId"`
	SyncState    string         `gorm:"default:none" json:"syncState"`
	Sort         int            `json:"sort"`
}

type Scene struct {
	Base
	CampaignID   uuid.UUID      `gorm:"type:uuid;index;not null" json:"campaignId"`
	FolderID     *uuid.UUID     `gorm:"type:uuid;index" json:"folderId"`
	Name         string         `gorm:"not null" json:"name"`
	Summary      string         `json:"summary"`
	SceneStatus  string         `gorm:"default:draft" json:"sceneStatus"`
	MapAssetID   *uuid.UUID     `gorm:"type:uuid" json:"mapAssetId"`
	ImageURL     string         `gorm:"-" json:"imageUrl"`
	Grid         datatypes.JSON `gorm:"type:jsonb" json:"grid"`
	SystemData   datatypes.JSON `gorm:"type:jsonb" json:"systemData"`
	FoundryUUID  string         `json:"foundryUuid"`
	FoundryDocID string         `json:"foundryDocId"`
	SyncState    string         `gorm:"default:none" json:"syncState"`
	Sort         int            `json:"sort"`
}

type ImageEntry struct {
	Base
	CampaignID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"campaignId"`
	FolderID     *uuid.UUID `gorm:"type:uuid;index" json:"folderId"`
	Name         string     `gorm:"not null" json:"name"`
	AssetID      *uuid.UUID `gorm:"type:uuid" json:"assetId"`
	ImageURL     string     `gorm:"-" json:"imageUrl"`
	ImageAlign   string     `gorm:"default:center" json:"imageAlign"`
	PushAs       string     `gorm:"default:empty_actor" json:"pushAs"`
	Notes        string     `json:"notes"`
	FoundryUUID  string     `json:"foundryUuid"`
	FoundryDocID string     `json:"foundryDocId"`
	SyncState    string     `gorm:"default:none" json:"syncState"`
	Sort         int        `json:"sort"`
}

type Asset struct {
	Base
	CampaignID  uuid.UUID  `gorm:"type:uuid;index;not null" json:"campaignId"`
	Name        string     `json:"name"`
	Kind        string     `json:"kind"`
	Source      string     `json:"source"`
	StorageURL  string     `json:"storageUrl"`
	Mime        string     `json:"mime"`
	Width       int        `json:"width"`
	Height      int        `json:"height"`
	DalleJobID  *uuid.UUID `gorm:"type:uuid" json:"dalleJobId"`
	FoundryPath string     `json:"foundryPath"`
	SourceRef   string     `gorm:"index" json:"-"`
}

type Memory struct {
	Base
	CampaignID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"campaignId"`
	SubjectType  string     `gorm:"default:campaign" json:"subjectType"`
	SubjectID    *uuid.UUID `gorm:"type:uuid;index" json:"subjectId"`
	SubjectLabel string     `json:"subjectLabel"`
	Title        string     `gorm:"not null" json:"title"`
	Body         string     `json:"body"`
	Kind         string     `gorm:"default:note" json:"kind"`
	Level        string     `gorm:"default:info;index" json:"level"`
	Acknowledged bool       `gorm:"default:false" json:"acknowledged"`
	Pinned       bool       `json:"pinned"`
}

type Tag struct {
	Base
	CampaignID uuid.UUID `gorm:"type:uuid;index;not null" json:"campaignId"`
	Name       string    `json:"name"`
	Color      string    `json:"color"`
}

type EntityTag struct {
	Base
	CampaignID uuid.UUID `gorm:"type:uuid;index;not null" json:"campaignId"`
	TagID      uuid.UUID `gorm:"type:uuid;index;not null" json:"tagId"`
	EntityType string    `gorm:"index" json:"entityType"`
	EntityID   uuid.UUID `gorm:"type:uuid;index" json:"entityId"`
}

type Link struct {
	Base
	CampaignID uuid.UUID `gorm:"type:uuid;index;not null" json:"campaignId"`
	FromType   string    `gorm:"index" json:"fromType"`
	FromID     uuid.UUID `gorm:"type:uuid;index" json:"fromId"`
	ToType     string    `gorm:"index" json:"toType"`
	ToID       uuid.UUID `gorm:"type:uuid;index" json:"toId"`
	Kind       string    `gorm:"default:related" json:"kind"`
}

type DalleJob struct {
	Base
	CampaignID    uuid.UUID  `gorm:"type:uuid;index;not null" json:"campaignId"`
	Prompt        string     `json:"prompt"`
	Size          string     `json:"size"`
	Quality       string     `json:"quality"`
	Style         string     `json:"style"`
	Status        string     `gorm:"default:queued" json:"status"`
	ResultAssetID *uuid.UUID `gorm:"type:uuid" json:"resultAssetId"`
	Error         string     `json:"error"`
}

type FoundryConnection struct {
	Base
	CampaignID   uuid.UUID  `gorm:"type:uuid;index;not null" json:"campaignId"`
	Label        string     `json:"label"`
	PairingToken string     `json:"-"`
	Status       string     `gorm:"default:disconnected" json:"status"`
	LastSeenAt   *time.Time `json:"lastSeenAt"`
}
