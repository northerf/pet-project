# WebSocket Notifications

Этот документ описывает функциональность WebSocket для уведомлений в проекте.

## Обзор

Система уведомлений поддерживает как REST API, так и WebSocket соединения для real-time уведомлений.

## Архитектура

### Компоненты

1. **ClientManager** (`internal/realtime/realtime.go`) - управляет WebSocket соединениями
2. **NotificationWSHandler** (`internal/handler/notification_ws.go`) - обработчик WebSocket запросов
3. **NotificationService** (`internal/service/notification.go`) - сервис с интеграцией WebSocket

### Поток данных

1. Пользователь подключается к WebSocket через `/ws/notifications`
2. При создании уведомления через REST API, оно автоматически отправляется через WebSocket
3. Клиент получает уведомления в реальном времени

## API Endpoints

### REST API

- `POST /notification/` - создать уведомление
- `GET /notification/` - получить уведомления пользователя
- `POST /notification/mark-read` - отметить уведомления как прочитанные
- `GET /notification/unread-count` - получить количество непрочитанных уведомлений

### WebSocket

- `GET /ws/notifications` - WebSocket соединение для получения уведомлений

## Использование

### Подключение к WebSocket

```javascript
// Подключение к WebSocket
const ws = new WebSocket('ws://localhost:8080/ws/notifications');

ws.onopen = function(event) {
    console.log('Connected to WebSocket');
};

ws.onmessage = function(event) {
    const notification = JSON.parse(event.data);
    console.log('Received notification:', notification);
};

ws.onclose = function(event) {
    console.log('Disconnected from WebSocket');
};
```

### Создание уведомления

```bash
curl -X POST http://localhost:8080/notification/ \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "task_assigned",
    "message": "You have been assigned a new task"
  }'
```

### Получение уведомлений

```bash
curl -X GET "http://localhost:8080/notification/?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Структура уведомления

```json
{
  "id": 1,
  "user_id": 123,
  "type": "task_assigned",
  "message": "You have been assigned a new task",
  "is_read": false,
  "created_at": "2024-01-01T12:00:00Z"
}
```

## Типы уведомлений

- `task_assigned` - назначена задача
- `task_completed` - задача завершена
- `comment_added` - добавлен комментарий
- `project_updated` - проект обновлен
- `test` - тестовое уведомление

## Безопасность

- Все WebSocket соединения требуют JWT токен
- Токен проверяется middleware перед установкой соединения
- Пользователь получает только свои уведомления

## Тестирование

Для тестирования WebSocket функциональности используйте файл `examples/websocket_client.html`:

1. Запустите сервер: `go run cmd/pet-project/main.go`
2. Откройте `examples/websocket_client.html` в браузере
3. Введите JWT токен и нажмите "Connect"
4. Создайте уведомление через форму
5. Уведомление должно появиться в реальном времени

## Обработка ошибок

### WebSocket ошибки

- `401 Unauthorized` - неверный или отсутствующий токен
- `400 Bad Request` - неверный формат запроса
- Соединение закрывается при ошибках

### REST API ошибки

- `400 Bad Request` - неверные параметры
- `401 Unauthorized` - неавторизованный доступ
- `500 Internal Server Error` - внутренняя ошибка сервера

## Мониторинг

### Подключенные пользователи

```go
connectedUsers := clientManager.GetConnectedUsers()
fmt.Printf("Connected users: %v\n", connectedUsers)
```

### Статистика

- Количество активных соединений
- Количество отправленных уведомлений
- Ошибки соединений

## Производительность

- Поддержка множественных соединений
- Буферизация сообщений (канал на 10 сообщений)
- Автоматическое закрытие неактивных соединений
- Ping/Pong для поддержания соединений

## Расширение функциональности

### Добавление новых типов уведомлений

1. Добавьте новый тип в документацию
2. Используйте его в сервисах при создании уведомлений

### Broadcast уведомления

```go
// Отправить уведомление всем подключенным пользователям
clientManager.Broadcast(notification)
```

### Персональные уведомления

```go
// Отправить уведомление конкретному пользователю
clientManager.Send(userID, notification)
``` 