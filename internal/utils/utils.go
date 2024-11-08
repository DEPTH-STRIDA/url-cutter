package utils

import (
	"app/internal/cache"
	"app/internal/logger"
	"app/internal/models"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"gorm.io/gorm"
)

// GenerateQRCode генерирует QR-код для заданного URL и возвращает его как строку Base64
func GenerateQRCode(url string) (string, error) {
	// Генерация QR-кода
	qrCode, err := qr.Encode(url, qr.L, qr.Unicode)
	if err != nil {
		return "", err
	}

	// Установка размера QR-кода
	qrCode, err = barcode.Scale(qrCode, 200, 200) // Измените размер по мере необходимости
	if err != nil {
		return "", err
	}

	// Создание буфера для хранения изображения
	var buf bytes.Buffer
	err = png.Encode(&buf, qrCode)
	if err != nil {
		return "", err
	}

	// Кодирование изображения в Base64
	imgBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return imgBase64, nil
}

// GenerateShortCode генерирует уникальную последовательность из 6 символов
func GenerateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // набор символов
	const length = 6                                                                 // длина кода

	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	shortCode := make([]byte, length)
	for i := range shortCode {
		shortCode[i] = charset[rand.Intn(len(charset))] // случайный выбор символа
	}
	return string(shortCode) // преобразование массива байтов в строку
}

func CreateNewValue(longUrl string) (string, error) {
	for i := 0; i < 10; i++ {

		shortUrl := GenerateShortCode()

		err := models.InsertNewUrl(shortUrl, longUrl)
		// Проверка на ошибку
		if err != nil {
			// Проверка на уникальность
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				logger.Log.Info("Сгенерирова занятая последовательность: ", err)
				continue
			}
			// Другая ошибка
			return "", err
		}

		cache.AddUrl(longUrl, shortUrl)

		return shortUrl, err
	}
	return "", fmt.Errorf("не удалось сгенерировать случайную последовательность после 10 попыток")
}

// isValidURL проверяет валидность URL
func IsValidUrl(urlString string) bool {
	// Убираем пробелы в начале и конце строки
	urlString = strings.TrimSpace(urlString)

	// Если URL начинается с www, добавляем протокол для корректной обработки
	if strings.HasPrefix(urlString, "www.") {
		urlString = "http://" + urlString
	}

	// Проверяем, начинается ли строка с http:// или https://
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		// Проверяем на локальные адреса
		if !strings.HasPrefix(urlString, "localhost") && !strings.HasPrefix(urlString, "127.0.0.1") {
			return false
		}
	}

	// Парсим URL
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return false // Если произошла ошибка, возвращаем false
	}

	// Проверяем, что hostname не пустой
	if parsedUrl.Hostname() == "" {
		return false
	}

	return true
}
