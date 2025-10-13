package store

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID          uuid.UUID `db:"id" json:"id"`
	OwnerID     uuid.UUID `db:"owner_id" json:"owner_id"` // FK to users
	Name        string    `db:"name" json:"name"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description"`
	LogoURL     string    `db:"logo_url" json:"logo_url"`
	BannerURL   string    `db:"banner_url" json:"banner_url"`
	IsPublic    bool      `db:"is_public" json:"is_public"`
	Location    string    `db:"location" json:"location"` // optional
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
