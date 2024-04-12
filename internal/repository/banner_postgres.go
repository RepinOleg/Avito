package repository

import (
	"database/sql"
	"errors"
	"fmt"

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

func (r *BannerPostgres) GetBanner(tagID, featureID int64) (*model.BannerContent, error) {
	var banner model.BannerContent
	query := fmt.Sprintf("SELECT content_title, content_text, content_url, is_active from %s b"+
		" JOIN %s bt ON b.banner_id=bt.banner_id WHERE feature_id = ($1) AND tag_id=($2) LIMIT 1;", bannerTable, tagBannerTable)

	row := r.db.QueryRow(query, featureID, tagID)
	var isActive bool
	if err := row.Scan(&banner.Title, &banner.Text, &banner.URL, &isActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &response.NotFoundError{Message: "banner not found"}
		}
		return nil, err
	}

	return &banner, nil
}

func (r *BannerPostgres) GetAllBanners(tagID, featureID, limit, offset int64) ([]response.BannerResponse200, error) {
	var banners []response.BannerResponse200
	query := fmt.Sprintf("SELECT b.banner_id, b.feature_id, b.content_title, b.content_text, b.content_url, b.is_active, b.created_at, b.updated_at, ARRAY_AGG(bt.tag_id) AS tag_ids "+
		"FROM %s b LEFT JOIN %s bt ON b.banner_id = bt.banner_id "+
		"WHERE 1=1 GROUP BY b.banner_id", bannerTable, tagBannerTable)

	var args []interface{}

	if tagID != 0 {
		query += " AND bt.tag_id IN (?)"
		args = append(args, tagID)
	}

	if featureID != 0 {
		query += " AND b.feature_id = ?"
		args = append(args, featureID)
	}

	query += " ORDER BY b.banner_id LIMIT $1 OFFSET $2"

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
	query = fmt.Sprintf("DELETE FROM %s WHERE banner_id = $1", bannerTable)
	result, err = tx.Exec(query, id)
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return false, err
	}

	var deleted = false
	if rowsAffected > 0 {
		deleted = true
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
