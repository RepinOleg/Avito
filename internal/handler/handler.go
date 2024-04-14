package handler

import (
	"net/http"
	"strings"

	"github.com/RepinOleg/Banner_service/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

type Route struct {
	Name    string
	Method  string
	Pattern string
}

type Routes []Route

func NewRouter(handlers *Handler) *mux.Router {
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
		case "sign-up":
			handlerFunc = handlers.SignUp
		case "sign-in":
			handlerFunc = handlers.SignIn
		}

		var resHandler http.Handler = handlerFunc

		if route.Name != "sign-up" && route.Name != "sign-in" {
			resHandler = handlers.TokenValidationMiddleware(resHandler, route.Name)
		}

		resHandler = Logger(resHandler, route.Name)

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
	Route{
		"sign-up",
		strings.ToUpper("Post"),
		"/sign-up",
	},
	Route{
		"sign-in",
		strings.ToUpper("Post"),
		"/sign-in",
	},
}
