package service

import (
	"time"

	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/repository"
)

type CacheService struct {
	cache repository.Cache
}

func NewCacheService(cache repository.Cache) *CacheService {
	return &CacheService{cache: cache}
}

func (s *CacheService) Create(id int64, banner model.BannerBody, duration time.Duration) {
	s.cache.AddBanner(id, banner, duration)
}
func (s *CacheService) Get(tagID, featureID int64) (*model.BannerContent, error) {
	return s.cache.GetBanner(tagID, featureID)
}
