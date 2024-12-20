package http

import (
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	API      *mux.Router
	shutdown chan os.Signal
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal) *App {
	api := App{
		API:      mux.NewRouter(),
		shutdown: shutdown,
	}
	return &api
}

// ServeHTTP API
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.API.ServeHTTP(w, r)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// API constructs an http.Handler with all application routes defined.
func API(shutdown chan os.Signal) http.Handler {
	// =========================================================================
	// Construct the web app api which holds all routes as well as common Middleware.

	api := NewApp(shutdown)

	// =========================================================================
	// Construct and attach relevant handlers to web app api

	api.API.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	router := mux.NewRouter()

	router.PathPrefix("/").Handler(api.API)

	return router
}
