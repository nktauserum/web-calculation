#!/usr/bin/bash

ROUTE=localhost
PORT=:8080

echo -e "Тест на обработку лишних символов"
curl -s -i --location "${ROUTE}${PORT}/api/v1/calculate" --header 'Content-Type: application/json' --data '{"expression": "2+2*a"}' | grep -q "HTTP/1.1 500" && echo "PASSED." || echo "FAILED."

echo -e "Тест на деление на ноль"
curl -s -i --location "${ROUTE}${PORT}/api/v1/calculate" --header 'Content-Type: application/json' --data '{"expression": "2+2/0"}' | grep -q "HTTP/1.1 500" && echo "PASSED." || echo "FAILED."

echo -e "Тест на приоритет операций"
curl -s -i --location "${ROUTE}${PORT}/api/v1/calculate" --header 'Content-Type: application/json' --data '{"expression": "2+2*2"}' | grep -q '{"id":' && echo "PASSED." || echo "FAILED."

echo -e "Тест на приоритет операций в скобках"
curl -s -i --location "${ROUTE}${PORT}/api/v1/calculate" --header 'Content-Type: application/json' --data '{"expression": 9}' | grep -q "HTTP/1.1 500" && echo "PASSED." || echo "FAILED."

echo -e "Тест на несоответствующие скобки"
curl -s -i --location "${ROUTE}${PORT}/api/v1/calculate" --header 'Content-Type: application/json' --data '{"expression": "(2+2"}' | grep -q "HTTP/1.1 500" && echo "PASSED." || echo "FAILED."

echo -e "Тест на недостаточно операндов"
curl -s -i --location "${ROUTE}${PORT}/api/v1/calculate" --header 'Content-Type: application/json' --data '{"expression": "+"}' | grep -q "HTTP/1.1 500" && echo "PASSED." || echo "FAILED."

echo -e "Тест на неизвестный оператор"
curl -s -i --location "${ROUTE}${PORT}/api/v1/calculate" --header 'Content-Type: application/json' --data '{"expression": "2 $ 2"}' | grep -q "HTTP/1.1 500" && echo "PASSED." || echo "FAILED."
