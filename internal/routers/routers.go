package routers

import (
	"fmt"
	"github.com/RepinOleg/Banner_service/internal/handlers"
	"github.com/RepinOleg/Banner_service/internal/middleware"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middleware.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"BannerGet",
		strings.ToUpper("Get"),
		"/banner",
		handlers.BannerGet,
	},

	Route{
		"BannerIdDelete",
		strings.ToUpper("Delete"),
		"/banner/{id}",
		handlers.BannerIdDelete,
	},

	Route{
		"BannerIdPatch",
		strings.ToUpper("Patch"),
		"/banner/{id}",
		handlers.BannerIdPatch,
	},

	Route{
		"BannerPost",
		strings.ToUpper("Post"),
		"/banner",
		handlers.BannerPost,
	},

	Route{
		"UserBannerGet",
		strings.ToUpper("Get"),
		"/user_banner",
		handlers.UserBannerGet,
	},
}
