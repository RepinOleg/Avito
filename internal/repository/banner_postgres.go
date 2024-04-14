package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/response"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type BannerPostgres struct {
	db *sqlx.DB
}

func NewBannerPostgres(db *sqlx.DB) *BannerPostgres {
	return &BannerPostgres{db: db}
}

func (r *BannerPostgres) GetBanner(tagID, featureID int64) (*model.BannerContent, bool, error) {
	var banner model.BannerContent
	query := fmt.Sprintf("SELECT content_title, content_text, content_url, is_active from %s b"+
		" JOIN %s bt ON b.banner_id=bt.banner_id WHERE feature_id = ($1) AND tag_id=($2) LIMIT 1;", bannerTable, tagBannerTable)

	row := r.db.QueryRow(query, featureID, tagID)
	var isActive bool
	if err := row.Scan(&banner.Title, &banner.Text, &banner.URL, &isActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, &response.NotFoundError{Message: "banner not found"}
		}
		return nil, false, err
	}

	return &banner, isActive, nil
}

func (r *BannerPostgres) GetAllBanners(tagID, featureID, limit, offset int64) ([]response.BannerResponse200, error) {
	var banners []response.BannerResponse200
	query := fmt.Sprintf("SELECT b.banner_id, b.feature_id, b.content_title, b.content_text, b.content_url, b.is_active, b.created_at, b.updated_at, ARRAY_AGG(bt.tag_id) AS tag_ids "+
		"FROM %s b LEFT JOIN %s bt ON b.banner_id = bt.banner_id WHERE 1=1", bannerTable, tagBannerTable)

	var args []interface{}
	var amount int
	if tagID != 0 {
		query += " AND bt.tag_id = $1"
		args = append(args, tagID)
		amount++
	}

	if featureID != 0 {
		query += " AND b.feature_id = $" + strconv.Itoa(amount+1)
		args = append(args, featureID)
		amount++
	}

	query += " GROUP BY b.banner_id ORDER BY b.banner_id LIMIT $" + strconv.Itoa(amount+1) + " OFFSET $" + strconv.Itoa(amount+2)

	args = append(args, limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			banner response.BannerResponse200
			tagIDs []int64
		)
		err = rows.Scan(&banner.BannerID, &banner.FeatureID, &banner.Content.Title, &banner.Content.Text, &banner.Content.URL, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt, pq.Array(&tagIDs))
		if err != nil {
			return nil, err
		}
		banner.TagIDs = tagIDs
		banners = append(banners, banner)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *BannerPostgres) AddBanner(banner model.BannerBody) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (feature_id) VALUES ($1) ON CONFLICT DO NOTHING", featureTable)

	_, err = tx.Exec(query, banner.FeatureID)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("error inserting to database feature: %s", err.Error())
	}
	query = fmt.Sprintf("INSERT INTO %s (feature_id, content_title, content_text, content_url, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING banner_id", bannerTable)
	var bannerID int64
	err = tx.QueryRow(query, banner.FeatureID, banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.IsActive).Scan(&bannerID)
	if err != nil {
		_ = tx.Rollback()
		return 0, fmt.Errorf("error inserting to database banner: %s", err.Error())
	}

	for _, tag := range banner.TagIDs {
		query = fmt.Sprintf("INSERT INTO %s (tag_id) VALUES ($1) ON CONFLICT DO NOTHING", tagTable)
		_, err = tx.Exec(query, tag)
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("error inserting to database tag: %s", err.Error())
		}
		query = fmt.Sprintf("INSERT INTO %s (banner_id, tag_id) VALUES ($1, $2)", tagBannerTable)
		_, err = tx.Exec(query, bannerID, tag)
		if err != nil {
			_ = tx.Rollback()
			return 0, fmt.Errorf("error inserting to database tag: %s", err.Error())
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return bannerID, nil
}

func (r *BannerPostgres) DeleteBanner(id int64) (bool, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return false, err
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE banner_id = $1", tagBannerTable)
	_, err = tx.Exec(query, id)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE banner_id = $1", bannerTable)
	result, err := tx.Exec(query, id)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	deleted := rowsAffected > 0

	query = fmt.Sprintf("DELETE FROM %s WHERE tag_id IN (SELECT t.tag_id FROM %s t LEFT JOIN %s bt ON t.tag_id = bt.tag_id WHERE bt.tag_id IS NULL)", tagTable, tagTable, tagBannerTable)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE feature_id IN (SELECT f.feature_id FROM %s f LEFT JOIN %s b ON f.feature_id = b.feature_id WHERE b.feature_id IS NULL)", featureTable, featureTable, bannerTable)
	_, err = tx.Exec(query)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return deleted, nil
}

func (r *BannerPostgres) PatchBanner(id int64, banner model.BannerBody) (bool, error) {
	var (
		content = banner.Content
		updated bool
	)
	query := fmt.Sprintf("UPDATE %s SET content_title=$1, content_text=$2, content_url=$3, updated_at=CURRENT_TIMESTAMP WHERE banner_id = $4", bannerTable)
	result, err := r.db.Exec(query, content.Title, content.Text, content.URL, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected > 0 {
		updated = true
	}
	return updated, nil
}
