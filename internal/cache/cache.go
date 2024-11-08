package cache

import (
	"app/internal/models"
	"errors"
	"sync"
)

// UrlCache представляет собой двусвязный кеш для хранения URL
type UrlCache struct {
	mu          sync.RWMutex      // Мьютекс для безопасного доступа к кэшу
	longToShort map[string]string // Хранит длинный URL по короткому
	shortToLong map[string]string // Хранит короткий URL по длинному
}

// Cache — глобальная переменная для кэша
var Cache *UrlCache

// InitCache инициализирует двусвязный кеш
func InitCache() error {
	// Валидация конфигурации
	if models.Config.TTL <= 0 {
		return errors.New("некорректное значение TTL: должно быть больше 0")
	}
	if models.Config.HardMaxCacheSize <= 0 {
		return errors.New("некорректное значение HardMaxCacheSize: должно быть больше 0")
	}
	if models.Config.MaxEntrySizes <= 0 {
		return errors.New("некорректное значение MaxEntrySize: должно быть больше 0")
	}

	// Инициализация кеша
	Cache = &UrlCache{
		longToShort: make(map[string]string),
		shortToLong: make(map[string]string),
	}

	return nil
}

// AddUrl добавляет новую пару (длинный URL : короткий URL) в кеш
func AddUrl(longUrl string, shortUrl string) {
	Cache.mu.Lock()
	defer Cache.mu.Unlock()

	// Проверка на существование длинного URL
	if _, exists := Cache.longToShort[longUrl]; !exists {
		Cache.longToShort[longUrl] = shortUrl
	}

	// Проверка на существование короткого URL
	if _, exists := Cache.shortToLong[shortUrl]; !exists {
		Cache.shortToLong[shortUrl] = longUrl
	}
}

// GetLongUrl возвращает длинный URL по короткому из кеша и обновляет его время жизни
func GetLongUrl(shortUrl string) (string, bool) {
	Cache.mu.RLock()
	defer Cache.mu.RUnlock()

	longUrl, exists := Cache.shortToLong[shortUrl]
	if exists {
		// Здесь можно добавить логику для обновления времени жизни (если это применимо)
	}
	return longUrl, exists
}

// GetShortUrl возвращает короткий URL по длинному из кеша и обновляет его время жизни
func GetShortUrl(longUrl string) (string, bool) {
	Cache.mu.RLock()
	defer Cache.mu.RUnlock()

	shortUrl, exists := Cache.longToShort[longUrl]
	if exists {
		// Здесь можно добавить логику для обновления времени жизни (если это применимо)
	}
	return shortUrl, exists
}
