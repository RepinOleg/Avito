package handler

import (
	"encoding/json"
	"github.com/RepinOleg/Banner_service/internal/dbs"
	"github.com/RepinOleg/Banner_service/internal/memorycache"
	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/response"
	"github.com/gorilla/mux"
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

func (h *Handler) GetAllBanners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

}

func (h *Handler) DeleteBannerID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bannerIDStr := params["id"]
	bannerID, err := strconv.ParseInt(bannerIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := h.db.DeleteBanner(bannerID)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		response.HandleError(w, "banner not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchBannerID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bannerIDStr := params["id"]
	bannerID, err := strconv.ParseInt(bannerIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//token := r.Header.Get("token")
	banner, err := readBody(r)

	ok, err := h.db.PatchBanner(bannerID, banner)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		response.HandleError(w, "banner not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PostBanner(w http.ResponseWriter, r *http.Request) {
	banner, err := readBody(r)

	bannerID, err := h.db.AddBanner(banner)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response201 := response.ModelResponse201{BannerID: bannerID}
	jsonResponse, err := json.MarshalIndent(response201, "", "\t")
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	// Добавление в кэш
	h.cache.Set(bannerID, banner, time.Minute*5)

	//отправка ответа
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println(err)
	}
}

func (h *Handler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	tagIDStr := r.FormValue("tag_id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJson(w, "Wrong parameter tag_id", http.StatusBadRequest)
		return
	}

	featureIDStr := r.FormValue("feature_id")
	featureID, err := strconv.ParseInt(featureIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJson(w, "Wrong parameter feature_id", http.StatusBadRequest)
		return
	}

	lastVersionStr := r.FormValue("use_last_revision")
	lastVersion, err := strconv.ParseBool(lastVersionStr)
	if err != nil {
		response.HandleErrorJson(w, "Wrong parameter use_last_version", http.StatusBadRequest)
		return
	}
	var content []model.BannerContent
	if lastVersion {
		content, err = h.db.GetBanner(tagID, featureID)
		if err != nil {
			response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		content, err = h.cache.GetBanner(tagID, featureID)
		if err != nil {
			response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	jsonResponse, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println(err)
	}
}

func readBody(r *http.Request) (model.BannerBody, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return model.BannerBody{}, err
	}

	var banner model.BannerBody
	err = json.Unmarshal(body, &banner)
	if err != nil {
		return model.BannerBody{}, err
	}
	return banner, nil
}
