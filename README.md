# url-cutter

Основные функции:
Сокращение URL: Пользователи могут вводить длинные URL-адреса, которые приложение сокращает до более коротких ссылок.

Хранение и извлечение URL: Сокращенные URL-адреса и их соответствующие длинные версии хранятся в базе данных и кэше для быстрого доступа.

Генерация QR-кодов: Для каждого сокращенного URL-адреса генерируется QR-код, который можно скачать.

Ограничение запросов: Внедрено ограничение количества запросов от одного IP-адреса для предотвращения злоупотреблений.

Логирование: Приложение ведет логирование событий и ошибок как в консоль, так и в файл.

.env пример
# APP
APP_IP=localhost
APP_PORT=8000

# DataBase
DBUSER=url_cutter_user
DBPASS=228
DBHOST=localhost
DBPORT=5432
DBNAME=url_cutter_database

# Cache
TTL=10 # В минутах
HARD_MAX_CACHE_SIZE=450
MAX_ENTRY_SIZE=2048
SHARDS=256

Имеет полностью реализованный, адаптивный дизайн
![2024-11-08_17-26-12](https://github.com/user-attachments/assets/a6883da6-3595-481e-bf40-4a0aa9768437)
![2024-11-08_17-26-32](https://github.com/user-attachments/assets/b03b5305-9855-4c96-8487-cd0c9a24ccff)
![2024-11-08_17-26-53](https://github.com/user-attachments/assets/f3e0358e-b7c0-40cb-9b16-8fda053b51e1)
