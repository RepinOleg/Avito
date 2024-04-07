package dbs

import (
	"fmt"
	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetBanner(tagID, featureID int64) ([]model.BannerContent, error) {
	var banners []model.BannerContent
	rows, err := r.db.Query("select content_title, content_text, content_url from banner b JOIN banner_tag bt ON b.banner_id=bt.banner_id WHERE feature_id = ($1) AND tag_id=($2);", featureID, tagID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var banner model.BannerContent

		if err = rows.Scan(&banner.Title, &banner.Text, &banner.URL); err != nil {
			return nil, err
		}

		banners = append(banners, banner)
	}
	return banners, nil
}

func (r *Repository) AddBanner(banner model.BannerBody) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// Вставка в таблицу feature
	_, err = tx.Exec("INSERT INTO feature (feature_id) VALUES ($1) ON CONFLICT DO NOTHING", banner.FeatureID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error inserting to database feature: %s", err.Error())
	}
	// Вставка в таблицу banner с использованием RETURNING
	var bannerID int64
	err = tx.QueryRow("INSERT INTO banner (feature_id, content_title, content_text, content_url, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING banner_id",
		banner.FeatureID, banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.IsActive).Scan(&bannerID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error inserting to database banner: %s", err.Error())
	}

	// Вставка в таблицу Tag и BannerTag
	for _, tag := range banner.TagIDs {
		_, err = tx.Exec("INSERT INTO tag (tag_id) VALUES ($1) ON CONFLICT DO NOTHING", tag)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error inserting to database tag: %s", err.Error())
		}
		// Добавление в таблицу banner_tag
		_, err = tx.Exec("INSERT INTO banner_tag (banner_id, tag_id) VALUES ($1, $2)", bannerID, tag)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error inserting to database tag: %s", err.Error())
		}
	}

	// Фиксация изменений в БД
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return bannerID, nil
}
