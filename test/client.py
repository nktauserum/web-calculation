import requests
import json
import errors
import time

class Calculator:
    # Основной путь к API оркестратора. 
    endpoint = ""

    def __init__(self, endpoint: str):
        self.endpoint = endpoint

    def _request(self, path: str, body: dict, token=None):
        # Выполняет HTTP запрос к API оркестратора с заданными параметрами
        response = None

        if token is not None:
            # авторизируемся
            response = requests.post(
                url=self.endpoint+path,
                data=json.dumps(body),
                headers={
                    "Authorization": f"Bearer {token}"
                }
            )
        else:
            response = requests.post(
                url=self.endpoint+path,
                data=json.dumps(body)
            )

        if response.status_code == 400:
            raise errors.BadRequestException(response_body=response.text)
        elif response.status_code == 401:
            raise errors.UnauthorizedException(response_body=response.text)
        elif response.status_code == 404:
            raise errors.NotFoundException(response_body=response.text)
        elif response.status_code == 500:
            raise errors.InternalServerErrorException(response_body=response.text) 
        
        return response 
               
    def register(self, username: str, email: str, password: str) -> str:
        # Регистрирует нового пользователя и возвращает токен
        response = self._request(path="/auth/register", body={
            "username": username,
            "email": email,
            "password": password
        })  
             
        json_response = response.json()
        return json_response["token"]
    
    def login(self, username: str, password: str) -> str:
        # Выполняет вход пользователя и возвращает токен
        response = self._request(path="/auth/login", body={
            "username": username,
            "password": password
        })  
             
        json_response = response.json()
        return json_response["token"]

    def calculate(self, expression: str, token: str) -> float:
        # Отправляет выражение на вычисление и ожидает результат
        response = self._request(path="/calculate", body={
            "expression": expression
        }, token=token)

        json_response = response.json()
        expr_id = int(json_response["id"])

        while self._expression(expr_id, token) is None:
            time.sleep(0.1)

        return self._expression(expr_id, token)

    def _expression(self, id: int, token:str) -> float | None:
        # Получает результат вычисления выражения по его идентификатору
        response = self._request(path="/expressions/"+str(id), body=None, token=token)
        json_response = response.json()

        if json_response["status"]:
            return json_response["result"]
        else:
            return None