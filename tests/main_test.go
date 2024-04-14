package tests

import (
	"fmt"
	"github.com/RepinOleg/Banner_service/internal/service"
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
	repo    *repository.Repository
	service *service.Service
	token   string
}

func (s *APITestSuite) SetupSuite() {
	cfg := repository.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test",
		Password: "test",
		DB:       "postgres",
	}

	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)

	connect, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		s.FailNow("Failed to connect to postgres: ", err)
	}

	repo := repository.NewRepository(connect)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	s.handler = handlers
	s.service = services
	s.repo = repo

	if err = s.populateDB(); err != nil {
		s.FailNow("Failed to populate DB", err)
	}

	s.populateCache()
	s.token = s.createToken()
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
		_, err = s.repo.AddBanner(banner)
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

		s.repo.Cache.SetBanner(i, banner, time.Minute*5)
		active = !active
	}
}

func (s *APITestSuite) createToken() string {
	user := model.User{
		Username: "test2",
		Password: "test2",
	}
	_, err := s.service.Authorization.CreateUser(user)
	if err != nil {
		s.FailNow("register failed", err.Error())
	}
	token, err := s.service.Authorization.GenerateToken(user.Username, user.Password, "user")
	if err != nil {
		s.FailNow("Authorization failed", err.Error())
	}
	return token
}
