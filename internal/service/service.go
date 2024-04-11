package service

import (
	"time"

	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/repository"
	"github.com/RepinOleg/Banner_service/internal/response"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password, role string) (string, error)
	ParseToken(token string, adminFlag bool) (int, error)
}

type Banner interface {
	Create(banner model.BannerBody) (int64, error)
	GetAll(tagID, featureID, limit, offset int64) ([]response.BannerResponse200, error)
	Get(tagID, featureID int64) (*model.BannerContent, error)
	Delete(id int64) (bool, error)
	Update(id int64, newBanner model.BannerBody) (bool, error)
}

type Cache interface {
	Create(id int64, banner model.BannerBody, duration time.Duration)
	Get(tagID, featureID int64) (*model.BannerContent, error)
}

type Service struct {
	Authorization
	Banner
	Cache
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Banner:        NewBannerService(repos.Banner),
		Cache:         NewCacheService(repos.Cache),
	}
}
