package web

import (
	"app/internal/logger"
	"app/internal/models"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

// WebApp веб приложение. Мозг программы, который использует большинство других приложений
type WebApp struct {
	Router        *mux.Router                   // Маршрутизатор
	TemplateCache map[string]*template.Template // Карта шаблонов
}

// NewWebApp создает и возвращает веб приложение
func NewWebApp() (*WebApp, error) {
	// Загрузка шаблонов
	templateCache, err := NewTemplateCache("./ui/html/")
	if err != nil {
		return nil, err
	}

	app := WebApp{
		TemplateCache: templateCache,
	}
	// Установка параметров
	app.Router = app.SetRoutes()
	return &app, nil
}

// HandleUpdates запускает HTTP сервер
func (app *WebApp) HandleUpdates() error {
	logger.Log.Info("Запуск сервера по адрессу " + models.Config.WebAppConfig.AppIP + ":" + models.Config.WebAppConfig.AppPort)

	err := http.ListenAndServe(models.Config.WebAppConfig.AppIP+":"+models.Config.WebAppConfig.AppPort, app.Router)
	// err := http.ListenAndServe("localhost:"+models.Config.WebAppConfig.AppPort, app.Router)
	if err != nil {
		return fmt.Errorf("ошибка при запуске сервера: %v", err)
	}
	return nil
}
