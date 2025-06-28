# Go Service

**Go Service** — серверное приложение на Go для управления проектами, задачами, комментариями и уведомлениями с поддержкой real-time WebSocket-уведомлений и JWT-аутентификации.

---

## 🚀 Быстрый старт

1. **Клонируйте репозиторий:**

2. **Настройте базу данных (PostgreSQL):**
   - Проверьте строку подключения в `cmd/pet-project/main.go` (или используйте docker-compose).
   - Примените миграции из `init.sql`.

3. **Установите зависимости:**
   ```sh
   go mod tidy
   ```

4. **Запустите сервер:**
   ```sh
   go run cmd/pet-project/main.go
   ```

5. **(Для теста фронта) Запустите локальный сервер для примера:**
   ```sh
   cd examples
   python3 -m http.server 8081
   # Откройте http://localhost:8081/websocket_client.html
   ```

---

## 🛡️ Аутентификация

- Используется JWT.
- Для всех защищённых маршрутов требуется токен в заголовке:  
  `Authorization: Bearer <ваш_токен>`

---

## 📝 Основные функции и примеры запросов

### 1. Регистрация пользователя

```sh
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"yourpassword"}'
```

### 2. Авторизация (получение JWT)

```sh
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"yourpassword"}'
```
**Ответ:**  
`{"token": "ваш_jwt_токен"}`

---

### 3. Проекты

- **Создать проект**
  ```sh
  curl -X POST http://localhost:8080/projects/ \
    -H "Authorization: Bearer <ваш_токен>" \
    -H "Content-Type: application/json" \
    -d '{"name":"My Project","description":"Описание"}'
  ```

- **Получить проект**
  ```sh
  curl -X GET http://localhost:8080/projects/1 \
    -H "Authorization: Bearer <ваш_токен>"
  ```

- **Обновить проект**
  ```sh
  curl -X PUT http://localhost:8080/projects/1 \
    -H "Authorization: Bearer <ваш_токен>" \
    -H "Content-Type: application/json" \
    -d '{"name":"New Name","description":"Новое описание"}'
  ```

- **Удалить проект**
  ```sh
  curl -X DELETE http://localhost:8080/projects/1 \
    -H "Authorization: Bearer <ваш_токен>"
  ```

---

### 4. Задачи

- **Создать задачу**
  ```sh
  curl -X POST http://localhost:8080/tasks/ \
    -H "Authorization: Bearer <ваш_токен>" \
    -H "Content-Type: application/json" \
    -d '{"title":"Task 1","description":"Описание","status":"pending","priority":"medium","assignedTo":2,"projectID":1,"dueDate":"2024-07-01T12:00:00Z"}'
  ```

- **Получить задачу**
  ```sh
  curl -X GET http://localhost:8080/tasks/1 \
    -H "Authorization: Bearer <ваш_токен>"
  ```

- **Обновить задачу**
  ```sh
  curl -X PUT http://localhost:8080/tasks/1 \
    -H "Authorization: Bearer <ваш_токен>" \
    -H "Content-Type: application/json" \
    -d '{"taskID":1,"assignedTo":2,"title":"New Title","status":"in_progress","priority":"high","description":"Новое описание"}'
  ```

- **Удалить задачу**
  ```sh
  curl -X DELETE http://localhost:8080/tasks/1 \
    -H "Authorization: Bearer <ваш_токен>"
  ```

---

### 5. Комментарии

- **Добавить комментарий**
  ```sh
  curl -X POST http://localhost:8080/comments/ \
    -H "Authorization: Bearer <ваш_токен>" \
    -H "Content-Type: application/json" \
    -d '{"taskID":1,"userID":2,"text":"Комментарий"}'
  ```

- **Удалить комментарий**
  ```sh
  curl -X DELETE http://localhost:8080/comments/1 \
    -H "Authorization: Bearer <ваш_токен>"
  ```

- **Получить комментарии к задаче**
  ```sh
  curl -X GET http://localhost:8080/comments/task/1 \
    -H "Authorization: Bearer <ваш_токен>"
  ```

- **Получить комментарии пользователя**
  ```sh
  curl -X GET http://localhost:8080/comments/user/2 \
    -H "Authorization: Bearer <ваш_токен>"
  ```

- **Обновить текст комментария**
  ```sh
  curl -X PUT http://localhost:8080/comments/1 \
    -H "Authorization: Bearer <ваш_токен>" \
    -H "Content-Type: application/json" \
    -d '{"text":"Новый текст"}'
  ```

---

### 6. Уведомления (REST)

- **Создать уведомление**
  ```sh
  curl -X POST http://localhost:8080/notification/ \
    -H "Authorization: Bearer <ваш_токен>" \
    -H "Content-Type: application/json" \
    -d '{"type":"test","message":"Test notification"}'
  ```

- **Получить уведомления**
  ```sh
  curl -X GET "http://localhost:8080/notification/?limit=10&offset=0" \
    -H "Authorization: Bearer <ваш_токен>"
  ```

- **Отметить уведомления как прочитанные**
  ```sh
  curl -X POST http://localhost:8080/notification/mark-read \
    -H "Authorization: Bearer <ваш_токен>" \
    -H "Content-Type: application/json" \
    -d '{"notification_ids":[1,2,3]}'
  ```

- **Получить количество непрочитанных**
  ```sh
  curl -X GET http://localhost:8080/notification/unread-count \
    -H "Authorization: Bearer <ваш_токен>"
  ```

---

### 7. Уведомления (WebSocket)

- **Подключение к WebSocket**
  - Используйте пример из `examples/websocket_client.html` или:
    ```js
    const ws = new WebSocket('ws://localhost:8080/ws/notifications?token=ВАШ_JWT_ТОКЕН');
    ws.onmessage = (event) => console.log(JSON.parse(event.data));
    ```

---

## 🖥️ Архитектура

- **Go + PostgreSQL**
- JWT-аутентификация
- REST API для всех сущностей
- WebSocket для real-time уведомлений
- CORS для фронта
- Пример клиента — в папке `examples/`

---

## 📦 Структура репозитория

```
pet-project/
├── cmd/
│   └── pet-project/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── middleware/
│   ├── realtime/
│   ├── repository/
│   └── service/
├── pkg/
│   └── model/
├── examples/
│   └── websocket_client.html
├── init.sql
├── README.md
├── LICENSE
├── go.mod
├── go.sum
└── .gitignore
```

---

## 📄 Лицензия

GNU General Public License v3.0 

---

## 🤝 Контакты и вклад

- Pull requests и идеи приветствуются!
- Telegram: @northerf (https://t.me/northerf)

---

**Удачного использования!** 
