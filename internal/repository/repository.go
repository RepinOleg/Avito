package repository

import (
	"time"

	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/response"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type Banner interface {
	GetBanner(tagID, featureID int64) (*model.BannerContent, error)
	GetAllBanners(tagID, featureID, limit, offset int64) ([]response.BannerResponse200, error)
	AddBanner(banner model.BannerBody) (int64, error)
	DeleteBanner(id int64) (bool, error)
	PatchBanner(id int64, banner model.BannerBody) (bool, error)
}

type Cache interface {
	AddBanner(id int64, item model.BannerBody, duration time.Duration)
	GetBanner(tagID, featureID int64) (*model.BannerContent, error)
}

type Repository struct {
	Authorization
	Banner
	Cache
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Banner:        NewBannerPostgres(db),
	}
}
