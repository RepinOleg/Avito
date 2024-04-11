package service

import (
	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/repository"
	"github.com/RepinOleg/Banner_service/internal/response"
)

type BannerService struct {
	repo repository.Banner
}

func NewBannerService(repo repository.Banner) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) Create(banner model.BannerBody) (int64, error) {
	return s.repo.AddBanner(banner)
}
func (s *BannerService) GetAll(tagID, featureID, limit, offset int64) ([]response.BannerResponse200, error) {
	return s.repo.GetAllBanners(tagID, featureID, limit, offset)
}
func (s *BannerService) Get(tagID, featureID int64) (*model.BannerContent, error) {
	return s.repo.GetBanner(tagID, featureID)
}
func (s *BannerService) Delete(id int64) (bool, error) {
	return s.repo.DeleteBanner(id)
}
func (s *BannerService) Update(id int64, newBanner model.BannerBody) (bool, error) {
	return s.repo.PatchBanner(id, newBanner)
}
