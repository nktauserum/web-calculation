# Сервис для вычисления арифметических выражений

## Описание

Этот сервис позволяет пользователям отправлять арифметические выражения через POST-запросы на сервер и получать в ответ рассчитанный результат.

## API Endpoint

Сервис имеет один endpoint:

- URL: /api/v1/calculate
- Метод: POST
- Тело запроса:

  {
  "expression": "выражение, которое ввёл пользователь"
  }

- Успешный ответ:

  {
  "result": "результат выражения"
  }

  - Код статуса: 200

- Ответы с ошибками:

  - Некорректное выражение:

    {
    "error": "Expression is not valid"
    }

    - Код статуса: 422

  - Внутренняя ошибка сервера:

    {
    "error": "Internal server error"
    }

    - Код статуса: 500

## Примеры использования

### Успешное вычисление

curl --location 'http://localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
"expression": "2+2\*2"
}'

- Ответ:

  {
  "result": "6"
  }

### Некорректное выражение

curl --location 'http://localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
"expression": "2+2\*a"
}'

- Ответ:

  {
  "error": "Expression is not valid"
  }

### Внутренняя ошибка сервера

curl --location 'http://localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
"expression": ""
}'

- Ответ:

  {
  "error": "Internal server error"
  }

## Установка и запуск проекта

1. Клонируйте репозиторий:

   git clone https://github.com/yourusername/calculator-web-service.git
   cd calculator-web-service

2. Запустите проект:

   go run ./cmd/main.go

По умолчанию сервис будет доступен по адресу http://localhost:8080.
