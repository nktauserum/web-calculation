# Сервис тестирования

Тестирование реализовано с помощью Python-скрипта, который проверяет основной функционал системы.

### Алгоритм работы тестов

1. Создается экземпляр класса Calculator для взаимодействия с API
2. Генерируются случайные данные для тестового пользователя
3. Последовательно выполняются все группы тестов
4. Для каждого теста:
   - Выполняется тестовый сценарий
   - Проверяется полученный результат
   - Выводится статус прохождения теста
5. В конце каждой группы тестов выводится статистика

### Структура тестов

1. **Тесты регистрации:**
   - Проверка успешной регистрации нового пользователя
   - Проверка запрета регистрации с пустым именем пользователя
   - Проверка запрета повторной регистрации

2. **Тесты авторизации:**
   - Проверка успешной авторизации существующего пользователя
   - Проверка запрета авторизации с неверным паролем
   - Проверка запрета авторизации несуществующего пользователя

3. **Тесты вычислений:**
   - Проверка корректности простых вычислений
   - Проверка корректности сложных выражений со скобками
   - Проверка обработки некорректных выражений
   - Проверка запрета вычислений без авторизации


### Запуск тестов

``` bash
python test/testing.py
```

### Зависимости

Для работы тестов требуется:
- **Python 3**
- Запущенный оркестратор на порту 8080
- Запущенный агент

### Конфигурация

Адрес тестируемого сервиса задается в переменной ENDPOINT в файле testing.py:
``` python
ENDPOINT = "http://localhost:8080/api/v1"
```
