# API для сервиса сокращения ссылок
API предназначен для создания коротких ссылок, которые перенаправляют пользователя на оригинальный URL. Сервис упрощает работу с длинными ссылками, делает их удобными для передачи, хранения и использования в интерфейсах.


## Документация API


### Конечные точки
**GET /health**
- **Описание**: проверка работы сервиса
- **Успешный ответ** (200 OK):
```json
{
  "status":  "ok",
  "service": "url-api"
}
```

**POST /shorten**
- **Описание**: принимает длинную ссылку, возвращает короткий идентификатор
- **Тело запроса**:
```json
{
  "base_url": "https://example.com",
  "day_life": 3,
  "code_length" : 6
}
```
- **Успешный ответ** (200 OK):
```json
{
  "status_code": 200,
	"description": "ok",
	"data": {
		"uniq_code": "rttRYU",
		"views": 0,
		"finally_date": "2026-03-25 21:01:52"
	}
}
```

**GET /{short_id}**
- **Описание**: редиректит на оригинальную ссылку
- **Успешный ответ** (301 Moved Permanently)


**GET /stats/{short_id}**
- **Описание**: возвращает количество переходов по ссылке
- **Успешный ответ** (200 OK):
```json
{
  "status_code": 200,
  "description": "ok",
  "data": 63
}
```


## База данных


### Конфигурация базы данных
- **База данных**: PostgreSQL
- **Драйвер**: postgres


### Схема базы данных
**url_table**
```sql
CREATE TABLE IF NOT EXISTS url_table (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    uniq_code TEXT,
    url_base TEXT,
    views INTEGER,
    finally_date TEXT
);

CREATE INDEX idx_url_table_uniq_code
ON url_table (uniq_code);
```


## Переменные окружения
```env
# DEV или PROD
MODE='DEV'

# Строка подключения к БД
CONNECTION_STRING='postgres://url_user:url_password@postgres:5432/url_db?sslmode=disable'

# Строка для генерации идентификатора
SYMBOLS='QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890'
```


## Как запустить


### Скопировать репозиторий
```bach
git clone https://github.com/theforps/URL_shortener_Go_SERVER.git

cd URL_shortener_Go_SERVER

git checkout upgrade
```


### Локальный запуск
1. **Скопировать конфигурацию переменных окружения:**
```bach
cp .env.example .env
```

2. **Установить зависимости:**
```bach
go mod download
go mod verify
```

3. **Запустить сервис:**
```bach
# Запустить код
go run cmd/main.go

# Создать исполняемый файл и запустить
go build -o api cmd/main.go
./api
```

**Сервис будет запущен по ссылке:** `http://localhost:5050`


### Docker-compose
```bash
docker-compose up --build

docker-compose down
```

**Сервис будет запущен по ссылке:** `http://localhost:8080`


## Тесты


### Алгоритм скриптов
1. **Поднятие контейнера с тестовой базой данных PostgreSQL**
2. **Очищение кэша и запуск тестов**
3. **Остановка и удаление, чтобы не осталось лишних контейнеров**


### Скрипты 
1. **Для операционной системы Windows**
```bash
./test.bat
```

2. **Для операционной системы Linux**
```bash
make test
```