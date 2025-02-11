# URL-Cutter

Сервис для сокращения длинных URL-адресов с возможностью генерации QR-кодов.

> 🌐 **Рабочий пример:** [golang-developer.ru](https://golang-developer.ru/)

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

### Алгоритм генерации коротких ссылок

1. **Детерминированная генерация**
   - Одинаковые URL всегда получают одинаковый короткий код
   - Экономия места в базе данных
   - Предсказуемость результата для пользователей
   - Улучшенная производительность кэширования

2. **Формат короткого URL**
   - Длина: 6 символов
   - Набор символов: a-z, A-Z, 0-9 (62 символа)
   - Возможные комбинации: 62^6 ≈ 56.8 миллиардов уникальных ссылок

3. **Процесс генерации**
   ```go
   // Для одинаковых URL будет сгенерирован одинаковый короткий код
   const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
   const length = 6
   ```
   - Проверка наличия URL в кэше
   - При отсутствии - генерация нового кода
   - Сохранение связки URL-код в базе и кэше
   - Повторные запросы того же URL вернут тот же код

4. **Преимущества подхода**
   - Отсутствие дубликатов в базе данных
   - Быстрый поиск по кэшу для популярных URL
   - Консистентность ссылок при повторных запросах
   - Эффективное использование памяти

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
<p align="center">
  <img src="https://github.com/user-attachments/assets/dbdcb0ae-3ed5-4aec-8e61-d8eab55d420c" width="80%" alt="Десктопная версия с пустым полем ввода" />
</p>
<p align="center">
  <img src="https://github.com/user-attachments/assets/1d0a7f56-4f64-4fce-89f2-41698b2ada15" width="80%" alt="Десктопная версия с qr кодом" />
</p>

<p align="center">
  <img src="https://github.com/user-attachments/assets/7c09b4ec-c219-42fa-a1c1-3e3243b4ecb2" width="32%" alt="Мобильная версия с пустым полем ввода" />
  <img src="https://github.com/user-attachments/assets/45ac28e1-50f6-4b18-a598-8154d38203f9" width="32%" alt="Мобильная версия с qr кодом" />
</p>
