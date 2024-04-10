package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestGetUserBannerFromDB() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=2&feature_id=1&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()

	s.handler.GetUserBanner(w, req)

	r := s.Require()
	r.Equal(http.StatusOK, w.Result().StatusCode)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	r.NoError(err)

	r.Equal("some_title1", actualResponse["title"])
	r.Equal("some_text1", actualResponse["text"])
	r.Equal("some_url1", actualResponse["url"])

}

func (s *APITestSuite) TestGetUserBannerFromCache() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=2&feature_id=1&use_last_revision=false", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()

	s.handler.GetUserBanner(w, req)

	r := s.Require()
	r.Equal(http.StatusOK, w.Result().StatusCode)

	var actualResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	r.NoError(err)

	r.Equal("some_title1", actualResponse["title"])
	r.Equal("some_text1", actualResponse["text"])
	r.Equal("some_url1", actualResponse["url"])
}

func (s *APITestSuite) TestGetUserBannerNotFoundTag() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=0&feature_id=1&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()

	s.handler.GetUserBanner(w, req)

	r := s.Require()
	r.Equal(http.StatusNotFound, w.Result().StatusCode)
}

func (s *APITestSuite) TestGetUserBannerBadRequest() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=str&feature_id=str&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()

	s.handler.GetUserBanner(w, req)

	r := s.Require()
	r.Equal(http.StatusBadRequest, w.Result().StatusCode)
}

func (s *APITestSuite) TestGetUserBanner() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=str&feature_id=str&use_last_revision=true", nil)
	req.Header.Set("accept", "application/json")
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()

	s.handler.GetUserBanner(w, req)

	r := s.Require()
	r.Equal(http.StatusBadRequest, w.Result().StatusCode)
}
