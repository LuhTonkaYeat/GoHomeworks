## Распределенная система для получения информации о репозиториях GitHub

### Архитектура
- **Collector** (gRPC сервер, порт 50051)
- **Gateway** (REST сервер, порт 8080)

### Запуск с помощью Docker (рекомендуемый способ)

1. **Клонируйте репозиторий:**
```bash
git clone https://github.com/LuhTonkaYeat/GoHomeworks.git
cd GoHomeworks/hw2
```

2. **Запустите сервисы:**
```bash
docker-compose up --build
```

3. **Используйте API:**
Получить информацию о репозитории:
GET http://localhost:8080/repo?owner={owner}&repo={repo}

# Информация о репозитории Go
curl "http://localhost:8080/repo?owner=golang&repo=go"

# Swagger UI (открыть в браузере)
http://localhost:8080/swagger/index.html

4. **Остановка:**
```bash
docker-compose down
```