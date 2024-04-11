package tests

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/RepinOleg/Banner_service/internal/handler"
	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite

	handler *handler.Handler
	cache   *repository.Cache
	repo    *repository.Repository
}

func (s *APITestSuite) SetupSuite() {
	cfg := repository.DBConfig{
		Addr:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DB:       "postgres",
	}

	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)

	connect, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		s.FailNow("Failed to connect to postgres: ", err)
	}

	cache := repository.New(5*time.Minute, 10*time.Minute)
	repo := repository.NewRepository(connect)
	s.handler = handler.NewHandler(repo, cache)
	s.cache = cache
	s.repo = repo

	if err = s.populateDB(); err != nil {
		s.FailNow("Failed to populate DB", err)
	}

	s.populateCache()
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

func (s *APITestSuite) populateDB() error {
	var active = true
	for i := int64(1); i < 10; i++ {
		banner := model.BannerBody{
			TagIDs:    []int64{i * 2, i * 3, i * 4},
			FeatureID: i,
			Content: model.BannerContent{
				Title: "some_title" + strconv.FormatInt(i, 10),
				Text:  "some_text" + strconv.FormatInt(i, 10),
				URL:   "some_url" + strconv.FormatInt(i, 10),
			},
			IsActive:   active,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Expiration: 0,
		}
		_, err := s.repo.AddBanner(banner)
		if err != nil {
			return err
		}
		active = !active
	}
	return nil
}

func (s *APITestSuite) populateCache() {
	var active = true
	for i := int64(1); i < 10; i++ {
		banner := model.BannerBody{
			TagIDs:    []int64{i * 2, i * 3, i * 4},
			FeatureID: i,
			Content: model.BannerContent{
				Title: "some_title" + strconv.FormatInt(i, 10),
				Text:  "some_text" + strconv.FormatInt(i, 10),
				URL:   "some_url" + strconv.FormatInt(i, 10),
			},
			IsActive:   active,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Expiration: 0,
		}
		s.cache.SetBanner(i, banner, time.Minute*5)
		active = !active
	}
}
