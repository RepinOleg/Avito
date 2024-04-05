package handler

import (
	"encoding/json"
	"fmt"
	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading request body")
		return
	}
	var banner model.BannerBody
	err = json.Unmarshal(body, &banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}
	_, err = h.db.Exec("INSERT INTO feature (feature_id) VALUES ($1)", banner.FeatureID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error inserting to database feature %s", err.Error())
		return
	}

	_, err = h.db.Exec("INSERT INTO banner (feature_id, content_title, content_text, content_url, is_active) VALUES ($1, $2, $3, $4, $5)",
		banner.FeatureID, banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.IsActive)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error inserting to database banner %s", err.Error())
		return
	}

	for _, tag := range banner.TagIDs {
		err = h.AddTag(tag)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error inserting to database tag %s", err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UserBannerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) AddTag(tag int64) error {
	_, err := h.db.Exec("INSERT INTO tag (tag_id) VALUES ($1)", tag)
	return err
}
