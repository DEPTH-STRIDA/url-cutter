# URL-Cutter

Сервис для сокращения длинных URL-адресов с возможностью генерации QR-кодов.

## Технологии

### Frontend
- Чистый HTML5, CSS3 и JavaScript без использования фреймворков
- Go HTML Templates (`.tmpl`) для серверного рендеринга
- Адаптивный дизайн для мобильных и десктопных устройств
- Отсутствие внешних зависимостей на клиентской стороне

### Backend
- Go (Golang)
- PostgreSQL
- In-memory кэширование

## Основные возможности

- **Сокращение URL:** Преобразование длинных URL в короткие, удобные для использования ссылки
- **QR-коды:** Автоматическая генерация QR-кодов для каждой сокращенной ссылки
- **Кэширование:** Быстрый доступ к часто используемым ссылкам
- **Rate Limiting:** Защита от злоупотребления через ограничение запросов
- **Мониторинг:** Подробное логирование всех операций

## Структура проекта

```
internal/
├── cache/    - Реализация кэширования
├── logger/   - Система логирования
├── models/   - Модели данных и работа с БД
├── utils/    - Вспомогательные функции
└── web/      - Веб-сервер и маршрутизация

ui/
├── html/     - HTML шаблоны
└── static/   - Статические файлы (CSS, JS, изображения)
```

## Endpoints

### Веб-интерфейс

- `/` - Главная страница
- `/go-cut` - Страница создания короткой ссылки
- `/rules` - Страница с правилами использования

### API

- `POST /shorten` - Создание короткой ссылки
  - Request: `{ "url": "https://example.com/very/long/url" }`
  - Response: `{ "short_url": "http://domain/u/abc123" }`

- `GET /u/{shortUrl}` - Переход по короткой ссылке
  - Автоматически перенаправляет на оригинальный URL

## Конфигурация

Настройки приложения задаются через `.env` файл:

```env
# Настройки приложения
APP_IP=localhost
APP_PORT=8000

# База данных
DBUSER=url_cutter_user
DBPASS=****
DBHOST=localhost
DBPORT=5432
DBNAME=url_cutter_database

# Параметры кэша
TTL=10           # Время жизни кэша в минутах
HARD_MAX_CACHE_SIZE=450
MAX_ENTRY_SIZE=2048
SHARDS=256
```

## Технические особенности

### Производительность
- Шардированный кэш для оптимизации памяти
- Подготовленные SQL-запросы
- Пулинг соединений с базой данных

### Безопасность
- Валидация входящих URL
- Rate limiting по IP
- Защита от SQL-инъекций

### Надёжность
- Многоуровневое логирование
- Graceful shutdown
- Обработка краевых случаев

## Запуск

1. Создайте `.env` файл на основе примера выше
2. Запустите PostgreSQL
3. Выполните:
```bash
go run main.go
```

Или используйте скрипт `run.bat` для Windows.

## Сборка

### Windows
```bash
go build main.go
```

### Linux
```bash
# Сборка под Linux x64 из Windows
set GOOS=linux
set GOARCH=amd64
go build main.go

# Сборка под Linux x64 из Linux
GOOS=linux GOARCH=amd64 go build main.go
```

### ARM (Raspberry Pi)
```bash
# 32-bit ARM (armv6)
set GOOS=linux
set GOARCH=arm
set GOARM=6
go build main.go

# 64-bit ARM
set GOOS=linux
set GOARCH=arm64
go build main.go
```

Результатом сборки будет исполняемый файл `main` (для Linux) или `main.exe` (для Windows).

## Примеры использования


