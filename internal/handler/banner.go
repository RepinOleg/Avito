package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/response"
	"github.com/gorilla/mux"
)

func (h *Handler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	tagIDStr := r.FormValue("tag_id")
	tagID, err := strconv.ParseInt(tagIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJSON(w, "wrong parameter tag_id or not found", http.StatusBadRequest)
		return
	}

	featureIDStr := r.FormValue("feature_id")
	featureID, err := strconv.ParseInt(featureIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJSON(w, "Wrong parameter feature_id or not found", http.StatusBadRequest)
		return
	}

	lastVersionStr := r.FormValue("use_last_revision")
	if lastVersionStr == "" {
		lastVersionStr = "false"
	}
	lastVersion, err := strconv.ParseBool(lastVersionStr)
	if err != nil {
		response.HandleErrorJSON(w, "Wrong parameter use_last_version", http.StatusBadRequest)
		return
	}

	var (
		content  *model.BannerContent
		isActive bool
	)

	if lastVersion {
		content, isActive, err = h.services.Banner.Get(tagID, featureID)
	} else {
		content, isActive, err = h.services.Cache.Get(tagID, featureID)
	}
	if err != nil {
		response.HandleError(w, err)
		return
	}

	header := r.Header.Get("token")
	if header == "" {
		http.Error(w, "пользователь не авторизован", http.StatusUnauthorized)
		return
	}
	if !isActive {
		headerParts := strings.Split(header, " ")
		_, err = h.services.Authorization.ParseToken(headerParts[1], true)
		if err != nil {
			response.HandleError(w, err)
			return
		}
	}

	if err != nil {
		response.HandleError(w, err)
		return
	}

	jsonResponse, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
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

	tagID, err := getOptionalInt64(r.FormValue("tag_id"), 0)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	featureID, err := getOptionalInt64(r.FormValue("feature_id"), 0)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := getOptionalInt64(r.FormValue("limit"), 10)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := getOptionalInt64(r.FormValue("offset"), 0)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	banners, err := h.services.Banner.GetAll(tagID, featureID, limit, offset)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(banners) == 0 {
		http.Error(w, "баннер не найден", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.MarshalIndent(banners, "", "\t")
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
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
	params := mux.Vars(r)
	bannerIDStr := params["id"]
	bannerID, err := strconv.ParseInt(bannerIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := h.services.Banner.Delete(bannerID)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
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
	params := mux.Vars(r)
	bannerIDStr := params["id"]
	bannerID, err := strconv.ParseInt(bannerIDStr, 10, 64)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	banner, err := readBody(r)
	if err != nil {
		response.HandleErrorJSON(w, "Неккоректные данные", http.StatusBadRequest)
		return
	}

	ok, err := h.services.Banner.Update(bannerID, banner)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
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

	banner, err := readBody(r)
	if err != nil {
		response.HandleErrorJSON(w, "Неккоректные данные", http.StatusBadRequest)
		return
	}

	bannerID, err := h.services.Banner.Create(banner)
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response201 := response.BannerResponse201{BannerID: bannerID}
	jsonResponse, err := json.MarshalIndent(response201, "", "\t")
	if err != nil {
		response.HandleErrorJSON(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	h.services.Cache.Create(bannerID, banner, time.Minute*5)

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
