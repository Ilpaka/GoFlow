1. Структура проекта
messenger-backend/
├── cmd/
│   └── app/
│       └── main.go
├── configs/
│   └── local.yaml
├── internal/
│   ├── app/
│   │   ├── app.go
│   │   └── container.go
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── user.go
│   │   ├── chat.go
│   │   ├── message.go
│   │   ├── session.go
│   │   └── common.go
│   ├── dto/
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── chat.go
│   │   └── message.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── user_repository.go
│   │   │   ├── chat_repository.go
│   │   │   ├── message_repository.go
│   │   │   └── session_repository.go
│   │   └── redis/
│   │       ├── presence_repository.go
│   │       ├── typing_repository.go
│   │       └── pubsub_repository.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── chat_service.go
│   │   ├── message_service.go
│   │   ├── presence_service.go
│   │   └── ws_service.go
│   ├── transport/
│   │   ├── http/
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── logging.go
│   │   │   │   └── recovery.go
│   │   │   ├── handler/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── user_handler.go
│   │   │   │   ├── chat_handler.go
│   │   │   │   └── message_handler.go
│   │   │   └── router.go
│   │   └── ws/
│   │       ├── handler.go
│   │       ├── hub.go
│   │       ├── client.go
│   │       ├── events.go
│   │       └── broadcaster.go
│   ├── pkg/
│   │   ├── auth/
│   │   │   ├── jwt.go
│   │   │   └── password.go
│   │   ├── logger/
│   │   │   └── logger.go
│   │   ├── validator/
│   │   │   └── validator.go
│   │   ├── errors/
│   │   │   └── errors.go
│   │   └── response/
│   │       └── response.go
│   └── migration/
│       ├── 001_init_users.sql
│       ├── 002_init_chats.sql
│       ├── 003_init_messages.sql
│       └── 004_init_sessions.sql
├── deployments/
│   ├── docker-compose.yml
│   └── Dockerfile
├── .env.example
├── go.mod
├── go.sum
└── README.md
2. Что делает каждый файл
cmd/app/main.go

Точка входа.
Что делает:

читает конфиг
создает приложение
запускает HTTP сервер
корректно завершает сервер по signal shutdown
internal/app/app.go

Сборка приложения.
Что делает:

хранит ссылки на config, logger, DB, Redis, router, ws hub
содержит Run() и Shutdown()
internal/app/container.go

Dependency injection вручную.
Что делает:

создаёт репозитории
создаёт сервисы
передаёт зависимости в handlers и WS layer
internal/config/config.go

Работа с конфигом.
Что делает:

описывает структуру конфига
читает YAML/env
валидирует обязательные поля
internal/domain/*.go

Чистые доменные сущности без HTTP и SQL.
Что делает:

описывает User, Chat, ChatMember, Message, RefreshSession
хранит enum-подобные типы: ChatType, MessageType, ChatRole
internal/dto/*.go

DTO для transport layer.
Что делает:

request/response модели
separate struct для JSON
не смешивает domain и HTTP payload
internal/repository/postgres/*.go

Работа с PostgreSQL.
Что делает:

CRUD и query-методы для users/chats/messages/sessions
никакой бизнес-логики
только чтение/запись в БД
internal/repository/redis/*.go

Работа с Redis.
Что делает:

online/offline presence
typing indicators
pub/sub для WS broadcast между инстансами
internal/service/*.go

Основная бизнес-логика.
Что делает:

auth_service: register/login/refresh/logout
user_service: me/update/search
chat_service: create direct/group, validate members
message_service: send/edit/delete/get history/read state
presence_service: online/offline/typing
ws_service: orchestration realtime-событий
internal/transport/http/middleware/*.go

HTTP middleware.
Что делает:

auth middleware: достаёт JWT, кладёт user id в context
logging middleware
recovery middleware
internal/transport/http/handler/*.go

HTTP handlers.
Что делает:

принимает запрос
валидирует DTO
вызывает service
возвращает JSON response
internal/transport/http/router.go

Регистрация маршрутов.
Что делает:

создаёт роуты /auth, /users, /chats, /messages, /ws
вешает middleware
internal/transport/ws/*.go

Realtime слой.
Что делает:

handler.go — upgrade HTTP -> WS
hub.go — реестр подключений
client.go — одна WS-сессия пользователя
events.go — формат входящих/исходящих WS-событий
broadcaster.go — рассылка в чат/пользователю
internal/pkg/auth/jwt.go

JWT utility.
Что делает:

generate access token
parse/validate access token
claims structure
internal/pkg/auth/password.go

Работа с паролем.
Что делает:

hash password
compare hash/password
internal/pkg/logger/logger.go

Инициализация логгера.

internal/pkg/validator/validator.go

Общая валидация DTO.

internal/pkg/errors/errors.go

App errors:

unauthorized
forbidden
not found
validation failed
conflict
internal/pkg/response/response.go

Унифицированный JSON response helper.

internal/migration/*.sql

SQL-миграции для таблиц.

deployments/docker-compose.yml

Поднимает:

app
postgres
redis
deployments/Dockerfile

Сборка Go сервиса в контейнер.

3. Что реализуем в MVP

Вот минимальный scope, который реально стоит делать первым:

Auth
Users
Direct chats
Group chats
Messages
Read state
WebSocket realtime
Presence + typing

Без:

файлов
реакций
голосовых
поиска по сообщениям
пушей
4. Порядок разработки

Вот правильная последовательность, чтобы проект не развалился:

каркас проекта
конфиг + логгер + Docker
БД + миграции
domain models
repositories
auth
users
chats
messages
websocket
presence/typing/read state
тесты и hardening