package main

import (
	cache "app/internal/cache"
	"app/internal/logger"
	"app/internal/models"
	"app/internal/web"
)

func main() {
	handleErr(logger.IniLogger())
	handleErr(models.InitConfig())
	handleErr(models.InitDataBase())
	handleErr(cache.InitCache())

	webApp, err := web.NewWebApp()
	handleErr(err)

	handleErr(webApp.HandleUpdates())
}

// handleErr проверяет наличие ошибки. В случае наличия ошикби, логгирует ее и останавливает программу.
func handleErr(err error) {
	if err != nil {
		logger.Log.Error("Ошибка при запуске приложения: %s", err)
		panic("")
	}
}
