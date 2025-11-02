# Wallet Service

Сервис для работы с кошельками пользователей в конкурентной среде. Использует горутины и мьютексы для обеспечения атомарности операций.

## Переменные окружения

Для инициализации установите следующие переменные:

- SERVER_PORT=8080
- DB_HOST=localhost
- DB_PORT=5432
- DB_USER=postgres
- DB_PASSWORD=password
- DB_NAME=wallet_db
- DB_SSLMODE=disable
