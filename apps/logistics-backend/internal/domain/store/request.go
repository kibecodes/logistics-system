package store

import "github.com/google/uuid"

type CreateStoreRequest struct {
	OwnerID     uuid.UUID `db:"admin_id" binding:"required"`
	Name        string    `db:"name" binding:"required"` // example:"Kevin's Electronics"
	Slug        string    `db:"slug" binding:"required"`
	LogoURL     string    `json:"logo_url"`                     // example:"https://cdn.fastabiz.com/logos/kevins.png"
	BannerURL   string    `json:"banner_url"`                   // example:"https://cdn.fastabiz.com/banners/kevins-banner.png"
	Description string    `db:"description" binding:"required"` // example:"Best electronics and accessories in Nairobi."
	IsPublic    bool      `db:"is_public" binding:"required"`
}

type UpdateStoreRequest struct {
	Column string      `json:"column" binding:"required"`
	Value  interface{} `json:"value" binding:"required"`
}

func (r *CreateStoreRequest) ToStore() *Store {
	return &Store{
		OwnerID:     r.OwnerID,
		Name:        r.Name,
		Slug:        r.Slug,
		LogoURL:     r.LogoURL,
		BannerURL:   r.BannerURL,
		Description: r.Description,
		IsPublic:    false,
	}
}
