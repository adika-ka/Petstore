# Реализация handlers Petstore

## Проверка логина пользователя

Я вручную проверил функциональность **регистрации и логина**, выполнив 3 шага:

- Зарегистрировал нового пользователя.
- Убедился, что пароль захеширован в базе.
- Успешно авторизовался с теми же данными и получил JWT-токен.

Код **не менялся**, проверку выполнял **3 раза** — результат стабильно положительный.

---

### Шаг 1. Регистрация пользователя

```bash
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "firstName": "Test",
    "lastName": "User",
    "email": "testuser@example.com",
    "password": "mypassword",
    "phone": "1234567890",
    "userStatus": 1
  }'
```
Ответ:
```bash
{
  "id": 1,
  "username": "testuser",
  "firstName": "Test",
  "lastName": "User",
  "email": "testuser@example.com",
  "password": "$2a$10$ohh3DrFPzn1pR1IGITqCWeEci0io7L4dZEzVLcRChouUCBagVY.1.",
  "phone": "1234567890",
  "userStatus": 1
}
```

### Шаг 2. Проверка пароля в базе

```bash
$ psql -h localhost -p 5432 -U postgres -d postgres
Password for user postgres: 
psql (16.8 (Ubuntu 16.8-0ubuntu0.24.04.1), server 15.12)
Type "help" for help.

postgres=# SELECT username, password FROM users WHERE username = 'testuser';
```
Ответ:
```bash
 username |                           password                           
----------+--------------------------------------------------------------
 testuser | $2a$10$ohh3DrFPzn1pR1IGITqCWeEci0io7L4dZEzVLcRChouUCBagVY.1.
(1 row)
```

### Шаг 3. Вход с теми же данными

```bash
$ curl "http://localhost:8080/user/login?username=testuser&password=mypassword"
```
Ответ:
```bash
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIn0.iOFkz34bQWVHEc6akqEKBi_vg9MXHZVvhH6pGrF8_rY"
}
```
