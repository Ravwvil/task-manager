# Task Manager API

REST API для управления I/O bound задачами.

## Установка и запуск

### Предварительные требования

- Docker установлен на вашей системе
- Порт 8080 должен быть свободен

### Клонирование репозитория

```bash
git clone https://github.com/Ravwvil/task-manager
cd task-manager
```

### Сборка и запуск

1. **Сборка Docker образа (с автоматическим запуском тестов):**
   ```bash
   docker build -t task-manager .
   ```
   При сборке автоматически запускаются тесты из `handlers_test.go`

2. **Запуск контейнера:**
   ```bash
   docker run --name task-manager -p 8080:8080 task-manager
   ```

### Тестирование

Тесты запускаются автоматически при сборке Docker образа. Если тесты не пройдут, сборка остановится.

**Просмотр подробных логов тестов:**
```bash
# Сборка с подробными логами тестов
docker build --progress=plain -t task-manager .
```

## API Endpoints

Сервер запускается на порту 8080. Все API endpoints имеют префикс `/api/v1`.

### 1. Проверка здоровья сервиса

```http
GET /health
```

---

### 2. Создание задачи

```http
POST /api/v1/tasks
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "description": "string" // обязательно
}
```

---

### 3. Получение списка всех задач

```http
GET /api/v1/tasks
```

---

### 4. Получение конкретной задачи

```http
GET /api/v1/tasks/{id}
```

**Параметры URL:**
- `id` - UUID задачи

---

### 5. Удаление задачи

```http
DELETE /api/v1/tasks/{id}
```
- `id` - UUID задачи

---
