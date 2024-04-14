package tests

import (
	"encoding/json"
	"github.com/RepinOleg/Banner_service/internal/handler"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestGetBannerFromDB() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=2&feature_id=1&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "Bearer "+s.token)
	w := httptest.NewRecorder()

	router := handler.NewRouter(s.handler)
	router.ServeHTTP(w, req)

	r := s.Require()
	r.Equal(http.StatusOK, w.Result().StatusCode)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	r.NoError(err)

	r.Equal("some_title1", actualResponse["title"])
	r.Equal("some_text1", actualResponse["text"])
	r.Equal("some_url1", actualResponse["url"])

}

func (s *APITestSuite) TestGetBannerFromCache() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=2&feature_id=1&use_last_revision=false", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "Bearer "+s.token)
	w := httptest.NewRecorder()
	router := handler.NewRouter(s.handler)
	router.ServeHTTP(w, req)

	r := s.Require()
	r.Equal(http.StatusOK, w.Result().StatusCode)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	r.NoError(err)

	r.Equal("some_title1", actualResponse["title"])
	r.Equal("some_text1", actualResponse["text"])
	r.Equal("some_url1", actualResponse["url"])
}

func (s *APITestSuite) TestGetBannerNotFound() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=0&feature_id=1&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "Bearer "+s.token)
	w := httptest.NewRecorder()
	router := handler.NewRouter(s.handler)
	router.ServeHTTP(w, req)

	r := s.Require()
	r.Equal(http.StatusNotFound, w.Result().StatusCode)
}

func (s *APITestSuite) TestGetBannerBadRequest() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=str&feature_id=str&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "Bearer "+s.token)
	w := httptest.NewRecorder()

	router := handler.NewRouter(s.handler)
	router.ServeHTTP(w, req)

	r := s.Require()
	r.Equal(http.StatusBadRequest, w.Result().StatusCode)
}

func (s *APITestSuite) TestGetBannerNotAuthorized() {
	req := httptest.NewRequest("GET", "/user_banner?tag_id=2&feature_id=1&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "WrongToken")

	w := httptest.NewRecorder()

	router := handler.NewRouter(s.handler)
	router.ServeHTTP(w, req)

	r := s.Require()
	r.Equal(http.StatusUnauthorized, w.Result().StatusCode)
}

func (s *APITestSuite) TestGetBannerNoAccess() {
	req := httptest.NewRequest("GET", "/user_banner?tag_id=4&feature_id=2&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "Bearer "+s.token)

	w := httptest.NewRecorder()

	router := handler.NewRouter(s.handler)
	router.ServeHTTP(w, req)

	r := s.Require()
	r.Equal(http.StatusForbidden, w.Result().StatusCode)
}
