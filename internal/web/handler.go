package web

import (
	cache "app/internal/cache"
	"app/internal/logger"
	"app/internal/models"
	"app/internal/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (wa *WebApp) HandleGoCut(w http.ResponseWriter, r *http.Request) {
	err := wa.render(w, "index.tmpl", nil)
	if err != nil {
		logger.Log.Error("Не удалось выполнить рендер: " + err.Error())
		http.Error(w, "Не удалось выполнить рендер: "+err.Error(), http.StatusInternalServerError)
	}
}

func (wa *WebApp) HandleMain(w http.ResponseWriter, r *http.Request) {
	err := wa.render(w, "go-cut.page.tmpl", nil)
	if err != nil {
		logger.Log.Error("Не удалось выполнить рендер: " + err.Error())
		http.Error(w, "Не удалось выполнить рендер: "+err.Error(), http.StatusInternalServerError)
	}
}

func (wa *WebApp) HandleRule(w http.ResponseWriter, r *http.Request) {
	err := wa.render(w, "rules.page.tmpl", nil)
	if err != nil {
		logger.Log.Error("Не удалось выполнить рендер: " + err.Error())
		http.Error(w, "Не удалось выполнить рендер: "+err.Error(), http.StatusInternalServerError)
	}
}

func (wa *WebApp) HandleShorten(w http.ResponseWriter, r *http.Request) {
	logger.Log.Info("Запрос от клиента IP --", r.URL.Path, " --/shorten")

	var jsonData map[string]interface{}

	// Декодирование JSON
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil && err != io.EOF {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		logger.Log.Error("Error parsing JSON: " + err.Error())
		return
	}
	defer r.Body.Close()

	// Проверка наличия long-url
	longUrl, ok := jsonData["long-url"].(string)
	if !ok {
		http.Error(w, "Данные не хранят long-url", http.StatusBadRequest)
		logger.Log.Error("Данные не хранят long-url")
		return
	}

	isValid := utils.IsValidUrl(longUrl)
	if !isValid {
		http.Error(w, "Неправильный url", http.StatusBadRequest)
		logger.Log.Error("От клиента поступил непрвильный url: ", longUrl)
		return
	}

	// КЕШ
	if shortUrl, ok := cache.GetShortUrl(longUrl); ok {
		respondWithShortUrl(w, shortUrl)
		return
	}

	// БД
	shortUrl, err := models.GetShortUrlByLong(longUrl)

	if err != nil {
		if errors.Is(err, models.ErrURLNotFound) {
			// Если URL не найден, генерируем новый
			shortUrl, err = utils.CreateNewValue(longUrl)
			if err != nil {
				http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
				logger.Log.Error("Ошибка при создании новой записи: ", err)
				return
			}
		} else {
			http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
			logger.Log.Error("Ошибка БД: ", err)
			return
		}
	}
	cache.AddUrl(longUrl, shortUrl)
	respondWithShortUrl(w, shortUrl)
}

func respondWithShortUrl(w http.ResponseWriter, shortUrl string) {
	shortUrl = "golang-developer.ru/u/" + shortUrl
	qrCode, err := utils.GenerateQRCode(shortUrl)
	if err != nil {
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		logger.Log.Error("Ошибка при генерации QR-кода: ", err)
		return
	}

	data := map[string]interface{}{"short-url": shortUrl, "img": qrCode}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		logger.Log.Error("Ошибка маршалинге JSON: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(dataJSON); err != nil {
		logger.Log.Error("Не удалось записать JSON в тело ответа: ", err)
	}
}

func (wa *WebApp) HandleGetLongUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortUrl := vars["shortUrl"] // Получаем короткий URL из маршрута

	logger.Log.Info("Получен короткий URL:", shortUrl)

	// КЕШ
	longUrl, ok := cache.GetLongUrl(shortUrl)
	if ok {
		http.Redirect(w, r, longUrl, http.StatusFound)
		return
	}

	// БД
	longUrl, err := models.GetLongUrlByShort(shortUrl)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		logger.Log.Error("Ошибка получения длинного URL: ", err)
		return
	}

	// Добавляем длинный URL в кеш
	cache.AddUrl(longUrl, shortUrl)
	http.Redirect(w, r, longUrl, http.StatusFound)
}

func (wa *WebApp) HandleYandexVerification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    </head>
    <body>Verification: f57292d76e700cf5</body>
</html>`))
}
