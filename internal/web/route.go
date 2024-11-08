package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Маршрутизатор
func (app *WebApp) SetRoutes() *mux.Router {
	router := mux.NewRouter()

	// Ограничение количества запросов от одного IP
	router.Use(LimitMiddleware)

	router.HandleFunc("/", app.HandleMain).Methods("GET")
	router.HandleFunc("/go-cut", app.HandleGoCut).Methods("GET")
	router.HandleFunc("/rules", app.HandleRule).Methods("GET")
	router.HandleFunc("/shorten", app.HandleShorten).Methods("POST")

	// Обработка запроса к длинному URL
	router.HandleFunc("/u/{shortUrl:.*}", app.HandleGetLongUrl).Methods("GET")

	staticDir := "./ui/static/"
	fileServer := http.FileServer(http.Dir(staticDir))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	return router
}
