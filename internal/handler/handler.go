package handler

import (
	"encoding/json"
	"github.com/RepinOleg/Banner_service/internal/dbs"
	"github.com/RepinOleg/Banner_service/internal/memorycache"
	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/response"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	db    *dbs.Repository
	cache *memorycache.Cache
}

func NewHandler(db *dbs.Repository, cache *memorycache.Cache) *Handler {
	return &Handler{
		db:    db,
		cache: cache,
	}
}

func (h *Handler) BannerGet(w http.ResponseWriter, r *http.Request) {
	tagIDStr := r.FormValue("tag_id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, "Wrong parameter tag_id", http.StatusBadRequest)
		return
	}

	featureIDStr := r.FormValue("feature_id")
	featureID, err := strconv.ParseInt(featureIDStr, 10, 64)
	if err != nil {
		response.HandleError(w, "Wrong parameter feature_id", http.StatusBadRequest)
		return
	}

	lastVersionStr := r.FormValue("use_last_version")
	lastVersion, err := strconv.ParseBool(lastVersionStr)
	if err != nil {
		response.HandleError(w, "Wrong parameter use_last_version", http.StatusBadRequest)
		return
	}
	var content []model.BannerContent
	if lastVersion {
		content, err = h.db.GetBanner(tagID, featureID)
		if err != nil {
			response.HandleError(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		content, err = h.cache.GetBanner(tagID, featureID)
		if err != nil {
			response.HandleError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	jsonResponse, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		response.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println(err)
	}

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

	bannerID, err := h.db.AddBanner(banner)
	if err != nil {
		response.HandleError(w, err.Error(), http.StatusInternalServerError)
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
