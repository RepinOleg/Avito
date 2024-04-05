package router

import (
	"github.com/RepinOleg/Banner_service/internal/handler"
	"github.com/RepinOleg/Banner_service/internal/middleware"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
}

type Routes []Route

func NewRouter(handlers *handler.Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handlerFunc http.HandlerFunc

		switch route.Name {
		case "BannerGet":
			handlerFunc = handlers.BannerGet
		case "BannerIdDelete":
			handlerFunc = handlers.BannerIDDelete
		case "BannerIdPatch":
			handlerFunc = handlers.BannerIDPatch
		case "BannerPost":
			handlerFunc = handlers.BannerPost
		case "UserBannerGet":
			handlerFunc = handlers.UserBannerGet
		}

		var resHandler http.Handler
		resHandler = handlerFunc
		resHandler = middleware.Logger(resHandler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(resHandler)
	}

	return router
}

var routes = Routes{

	Route{
		"BannerGet",
		strings.ToUpper("Get"),
		"/banner",
	},

	Route{
		"BannerIdDelete",
		strings.ToUpper("Delete"),
		"/banner/{id}",
	},

	Route{
		"BannerIdPatch",
		strings.ToUpper("Patch"),
		"/banner/{id}",
	},

	Route{
		"BannerPost",
		strings.ToUpper("Post"),
		"/banner",
	},

	Route{
		"UserBannerGet",
		strings.ToUpper("Get"),
		"/user_banner",
	},
}
