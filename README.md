# Сервис для вычисления арифметических выражений

## Описание

Этот сервис позволяет пользователям отправлять арифметические выражения через POST-запросы на сервер и получать в ответ рассчитанный результат.использовать

## Установка и запуск проекта

1. Клонируйте репозиторий:

   ```sh
   git clone https://github.com/veliashev/web-calculation

   cd web-calculation
   ```

2. Запустите проект:

   ```sh
   go run ./cmd/main.go
   ```

По умолчанию сервис будет доступен по адресу http://localhost:8080.

## API Endpoint

Сервис имеет один endpoint:

- URL: /api/v1/calculate
- Метод: POST
- Тело запроса:

  ```json
  { "expression": "ваше выражение" }
  ```

- Успешный ответ:

  ```json
  { "result": "результат выражения" }
  ```

  Код статуса: **200**

- Ответы с ошибками:

  - Некорректное выражение:

    ```json
    { "error": "некорректное выражение (подробности)..." }
    ```

    Код статуса: **422**

  - Внутренняя ошибка сервера:

    ```json
    { "error": "Internal server error" }
    ```

    Код статуса: **500**

## Примеры использования на порту :8080

Для запросов к серверу можно использовать консольную утилиту _curl_, можно графические утилиты, навроде _Postman_. Примеры приведены для _curl_. Используйте то, что вам удобнее.

По умолчанию используется порт **8080**, если у вас он занят, можете сменить его в файле [config.json](config/config.json).

### Успешное вычисление

```sh
curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "2+2\*2"}'
```

- Ответ:

  ```json
  { "result": 6 }
  ```

### Некорректное выражение

```sh
curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "2+2*a"}'
```

- Ответ:

  ```json
  { "error": "недопустимая цифра в выражении" }
  ```

### Внутренняя ошибка сервера

```sh
curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": 9}'
```

- Ответ:

  ```json
  {
    "error": "json: cannot unmarshal number into Go struct field Request.expression of type string"
  }
  ```
