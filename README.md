# http-gate-control

# HTTP Gate Control

## Установка

```bash
make build
docker run -e RELAY_URL=... -e OPEN_DELAY=... -e LOGGER_URL=...

Использование
Эндпоинт: GET /open
Swagger: GET /swagger/
Метрики: GET /metrics

