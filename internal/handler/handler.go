package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/RepinOleg/Banner_service/internal/memorycache"
	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/response"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	db    *sqlx.DB
	cache *memorycache.Cache
}

func NewHandler(db *sqlx.DB, cache *memorycache.Cache) *Handler {
	return &Handler{
		db:    db,
		cache: cache,
	}
}

func (h *Handler) BannerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

func (h *Handler) BannerIDDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) BannerIDPatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) BannerPost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var banner model.BannerBody
	err = json.Unmarshal(body, &banner)
	if err != nil {
		response.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		response.HandleError(w, "Error starting database transaction", http.StatusInternalServerError)
		return
	}

	// Вставка в таблицу feature
	_, err = tx.Exec("INSERT INTO feature (feature_id) VALUES ($1) ON CONFLICT DO NOTHING", banner.FeatureID)
	if err != nil {
		tx.Rollback()
		response.HandleError(w, fmt.Sprintf("Error inserting to database feature: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	// Вставка в таблицу banner с использованием RETURNING
	var bannerID int64
	err = tx.QueryRow("INSERT INTO banner (feature_id, content_title, content_text, content_url, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING banner_id",
		banner.FeatureID, banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.IsActive).Scan(&bannerID)
	if err != nil {
		tx.Rollback()
		response.HandleError(w, fmt.Sprintf("Error inserting to database banner: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Вставка в таблицу Tag и BannerTag
	for _, tag := range banner.TagIDs {
		err = h.AddTag(tx, tag)
		if err != nil {
			tx.Rollback()
			response.HandleError(w, fmt.Sprintf("Error inserting to database tag: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		// Добавление в таблицу banner_tag
		err = h.AddTagToBanner(tx, tag, bannerID)
		if err != nil {
			tx.Rollback()
			response.HandleError(w, fmt.Sprintf("Error inserting to database tag: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}

	// Фиксация изменений в БД
	if err := tx.Commit(); err != nil {
		response.HandleError(w, "Error committing database transaction", http.StatusInternalServerError)
		return
	}
	response201 := response.ModelResponse201{BannerID: bannerID}
	jsonResponse, err := json.MarshalIndent(response201, "", "\t")
	if err != nil {
		response.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	// Добавление в кэш
	h.cache.Set(bannerID, banner, time.Minute*5)

	//отправка ответа
	_, err = w.Write(jsonResponse)
	if err != nil {
		response.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UserBannerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddTag(tx *sql.Tx, tag int64) error {
	_, err := tx.Exec("INSERT INTO tag (tag_id) VALUES ($1) ON CONFLICT DO NOTHING", tag)
	return err
}

func (h *Handler) AddTagToBanner(tx *sql.Tx, tag, bannerID int64) error {
	_, err := tx.Exec("INSERT INTO banner_tag (banner_id, tag_id) VALUES ($1, $2)", bannerID, tag)
	return err
}
