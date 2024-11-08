package models

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DataBase *DataBaseStruct

type DataBaseStruct struct {
	*gorm.DB
}

// InitDataBase подключается к БД по данным из конфига. Устанавливает подключение в глобальную переменную.
func InitDataBase() error {

	// Стандартная строка для подключения к БД postgresql
	dbUri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		Config.DataBaseConfig.UserName,
		Config.DataBaseConfig.Password,
		Config.DataBaseConfig.Host,
		Config.DataBaseConfig.Port,
		Config.DataBaseConfig.DataBaseName)

	// Сохранение "Подключения" в переменную
	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		return err
	}
	DataBase = &DataBaseStruct{}

	DataBase.DB = conn

	// AutoMigrate создает таблицы, если они не существуют
	if err := DataBase.DB.AutoMigrate(&Urls{}); err != nil {
		return fmt.Errorf("не удалось выполнить миграцию: %w", err)
	}

	return nil
}

type Urls struct {
	gorm.Model
	ShortUrl string `gorm:"column:short_url;type:varchar(100);not null;unique"` // Короткая ссылка: до 100 символов, обязательное поле, уникальное
	LongUrl  string `gorm:"column:long_url;type:varchar(2000);not null;unique"` // Длинная ссылка: до 2000 символов, обязательное поле, уникальное
}

// TableName определяет имя таблицы для Group
func (Urls) TableName() string {
	return "urls"
}

// GetLongUrlByShort ищет длинный URL по короткому URL
func GetLongUrlByShort(shortUrl string) (string, error) {
	var url Urls
	result := DataBase.DB.Where("short_url = ?", shortUrl).First(&url)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", errors.New("URL not found")
	}

	return url.LongUrl, result.Error
}

var ErrURLNotFound = errors.New("URL not found")

// GetShortUrlByLong ищет короткий URL по длиному URL
func GetShortUrlByLong(longUrl string) (string, error) {
	var url Urls
	result := DataBase.DB.Where("long_url = ?", longUrl).First(&url)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", ErrURLNotFound
	}

	return url.ShortUrl, result.Error // Исправлено с LongUrl на ShortUrl
}

// InsertNewUrl вставляет новую запись в таблицу urls
func InsertNewUrl(shortUrl string, longUrl string) error {
	newUrl := Urls{
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}

	result := DataBase.DB.Create(&newUrl)
	return result.Error
}
