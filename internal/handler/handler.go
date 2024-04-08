package handler

import (
	"encoding/json"
	"fmt"
	"github.com/RepinOleg/Banner_service/internal/memorycache"
	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/repository"
	"github.com/RepinOleg/Banner_service/internal/response"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	db    *repository.Repository
	cache *memorycache.Cache
}

func NewHandler(db *repository.Repository, cache *memorycache.Cache) *Handler {
	return &Handler{
		db:    db,
		cache: cache,
	}
}

func (h *Handler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token == "" || token != "admin_token" && token != "user_token" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	tagIDStr := r.FormValue("tag_id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJson(w, "wrong parameter tag_id or not found", http.StatusBadRequest)
		return
	}

	featureIDStr := r.FormValue("feature_id")
	featureID, err := strconv.ParseInt(featureIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJson(w, "Wrong parameter feature_id or not found", http.StatusBadRequest)
		return
	}

	lastVersionStr := r.FormValue("use_last_revision")
	if lastVersionStr == "" {
		lastVersionStr = "false"
	}
	lastVersion, err := strconv.ParseBool(lastVersionStr)
	if err != nil {
		response.HandleErrorJson(w, "Wrong parameter use_last_version", http.StatusBadRequest)
		return
	}

	var content *model.BannerContent

	if lastVersion {
		content, err = h.db.GetBanner(tagID, featureID, token)
	} else {
		content, err = h.cache.GetBanner(tagID, featureID, token)
	}
	if err != nil {
		response.HandleError(w, err)
		return
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

func (h *Handler) GetAllBanners(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token == "" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	if token != "admin_token" {
		http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
		return
	}

	tagID, err := getOptionalInt64(r.FormValue("tag_id"), 0)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	featureID, err := getOptionalInt64(r.FormValue("feature_id"), 0)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := getOptionalInt64(r.FormValue("limit"), 10)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := getOptionalInt64(r.FormValue("offset"), 0)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	banners, err := h.db.GetAllBanners(tagID, featureID, limit, offset)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(banners) == 0 {
		http.Error(w, "баннер не найден", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.MarshalIndent(banners, "", "\t")
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

func (h *Handler) DeleteBannerID(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

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
		response.HandleError(w, &response.NotFoundError{Message: "banner not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PatchBannerID(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token != "admin_token" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}
	params := mux.Vars(r)
	bannerIDStr := params["id"]
	bannerID, err := strconv.ParseInt(bannerIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	banner, err := readBody(r)
	if err != nil {
		response.HandleErrorJson(w, "Неккоректные данные", http.StatusBadRequest)
		return
	}

	ok, err := h.db.PatchBanner(bannerID, banner)
	if err != nil {
		response.HandleErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		response.HandleError(w, &response.NotFoundError{Message: "banner not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PostBanner(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token == "" {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	if token != "admin_token" {
		http.Error(w, "Пользователь не имеет доступа", http.StatusForbidden)
		return
	}
	banner, err := readBody(r)
	if err != nil {
		response.HandleErrorJson(w, "Неккоректные данные", http.StatusBadRequest)
		return
	}

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

func getOptionalInt64(value string, defaultValue int64) (int64, error) {
	if value == "" {
		return defaultValue, nil
	}

	num, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return 0, err
	}

	if num < 0 {
		return 0, fmt.Errorf("value must be greater than or equal to 0, value = %s", value)
	}

	return num, nil
}
