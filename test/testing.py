from client import Calculator
from utils import pass_, fail, generate_random_string, bold, Counter, part

ENDPOINT = "http://localhost:8080/api/v1"

def registration_test():
    c = Counter()
    # 1: корректная регистрация
    try:
        c.all()
        calc.register(username=random_username, email=random_email, password="password123")
        pass_("Тест 1 пройден: Успешная регистрация")
        c.passed()
    except Exception as e:
        fail(f"Тест 1 не пройден: {e}")
        return

    # 2: пустое имя пользователя
    try:
        c.all()
        calc.register(username="", email=random_email, password="password123")
        fail("Тест 2 не пройден: Не должно допускаться пустое имя пользователя")
    except Exception as e:
        pass_("Тест 2 пройден: Пустое имя пользователя отклонено")
        c.passed()

    # 3: повторная регистрация
    try:
        c.all()
        calc.register(username=random_username, email=random_email, password="password123")
        fail("Тест 3 не пройден: Не должна допускаться повторная регистрация")
    except Exception as e:
        pass_("Тест 3 пройден: Повторная регистрация отклонена")
        c.passed()

    # Количество пройденных тестов к общему количеству
    print(f"Пройдено: ({c.final()[0]}/{c.final()[1]})")
    print()

def login_test():
    c = Counter()
    # 1: корректный вход
    try:
        c.all()
        token = calc.login(username=random_username, password="password123")
        pass_("Тест 1 пройден: Успешный вход")
        c.passed()
    except Exception as e:
        fail(f"Тест 1 не пройден: {e}")
        return

    # 2: неверный пароль
    try:
        c.all()
        calc.login(username=random_username, password="wrongpassword")
        fail("Тест 2 не пройден: Не должен допускаться вход с неверным паролем")
    except Exception as e:
        pass_("Тест 2 пройден: Вход с неверным паролем отклонен")
        c.passed()

    # 3: несуществующий пользователь
    try:
        c.all()
        calc.login(username="nonexistentuser", password="password123")
        fail("Тест 3 не пройден: Не должен допускаться вход несуществующего пользователя")
    except Exception as e:
        pass_("Тест 3 пройден: Вход несуществующего пользователя отклонен")
        c.passed()

    # Количество пройденных тестов к общему количеству
    print(f"Пройдено: ({c.final()[0]}/{c.final()[1]})")
    print()

def calculation_test():
    с_global = Counter()
    # 1: корректное вычисление
    try:
        с_global.all()
        token = calc.login(username=random_username, password="password123")
        result = calc.calculate("2+2*2", token)
        if float(result) == 6:
            с_global.passed()
            pass_("Тест 1 пройден: Успешное вычисление")
        else:
            fail(f"Тест 1 не пройден: Неверный результат {result}")
    except Exception as e:
        fail(f"Тест 1 не пройден: {e}")
        return    
    
    # 2: сложные выражения
    try:
        c_local = Counter()
        token = calc.login(username=random_username, password="password123")
        expressions = [
            ("(5+3)*2-4/2", 14),
            ("2+2*2+2/2", 7),
            ("(8+2*5)/(2+3)", 3.6),
            ("3*3/(2+1)-1", 2),
            ("10-2*3+4/2", 6)
        ]
        for expr, expected in expressions:
            с_global.all()
            c_local.all()
            result = float(calc.calculate(expr, token))
            if result != expected:
                fail(f"\tТест 2:")
                print(f"\tДля {expr} получено {result}, ожидалось {expected}")
            else:
                с_global.passed()
                c_local.passed()

        if c_local.final()[0] == c_local.final()[1]:
            pass_("Тест 2 пройден: Сложные выражения вычислены корректно")
        elif c_local.final()[0] == 0:
            fail("Тест 2 не пройден: сложные выражения не вычислены")
        else:
            part(f"Тест 2 пройден частично: ({c_local.final()[0]}/{c_local.final()[1]})")
            

    except Exception as e:
        fail(f"Тест 2 не пройден: {e}")

    # 3: некорректные выражения
    try:
        token = calc.login(username=random_username, password="password123")
        invalid_expressions = [
            "2++2",
            "2+2*",
            "(2+2",
            "2+2)",
            "2 $ 2",
            "2+2/0",
            "a+b*2"
        ]
        for expr in invalid_expressions:
            try:
                calc.calculate(expr, token)
                fail(f"Тест 3 не пройден: Выражение {expr} должно быть отклонено")
            except:
                continue
        pass_("Тест 3 пройден: Некорректные выражения отклонены")
    except Exception as e:
        pass_("Тест 3 пройден: Некорректные выражения отклонены")

    # 4: вычисление без авторизации
    try:
        calc.calculate("2+2", "invalid_token")
        fail("Тест 4 не пройден: Не должно допускаться вычисление без авторизации")
    except Exception as e:
        pass_("Тест 4 пройден: Вычисление без авторизации отклонено")

    print(f"Пройдено: ({с_global.final()[0]}/{с_global.final()[1]})")
    print()

if __name__ == "__main__":
    calc = Calculator(ENDPOINT)

    random_username = generate_random_string(8)
    random_email = f"{generate_random_string(8)}@{generate_random_string(6)}.com"
    
    bold("Регистрация:")
    registration_test()

    bold("Авторизация:")
    login_test()

    bold("Вычисление:")
    calculation_test()
