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
		case "GetAllBanners":
			handlerFunc = handlers.GetAllBanners
		case "DeleteBannerId":
			handlerFunc = handlers.DeleteBannerID
		case "PatchBannerId":
			handlerFunc = handlers.PatchBannerID
		case "PostBanner":
			handlerFunc = handlers.PostBanner
		case "GetUserBanner":
			handlerFunc = handlers.GetUserBanner
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
		"GetAllBanners",
		strings.ToUpper("Get"),
		"/banner",
	},

	Route{
		"DeleteBannerId",
		strings.ToUpper("Delete"),
		"/banner/{id}",
	},

	Route{
		"PatchBannerId",
		strings.ToUpper("Patch"),
		"/banner/{id}",
	},

	Route{
		"PostBanner",
		strings.ToUpper("Post"),
		"/banner",
	},

	Route{
		"GetUserBanner",
		strings.ToUpper("Get"),
		"/user_banner",
	},
}
