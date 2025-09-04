# 🛒 Marketplace - Микросервисная архитектура

Современная платформа электронной коммерции, построенная на микросервисной архитектуре с использованием Go, PostgreSQL и Apache Kafka.

## 🏗️ Архитектура

Проект состоит из 4 независимых микросервисов:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Service  │    │ Product Service │    │  Order Service  │    │Notification Svc│
│                 │    │                 │    │                 │    │                 │
│ • Регистрация   │    │ • Каталог       │    │ • Создание      │    │ • Email         │
│ • Авторизация   │    │ • Поиск         │    │ • Управление    │    │ • SMS           │
│ • Профили       │    │ • CRUD          │    │ • Статусы       │    │ • Push          │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │                       │
         └───────────────────────┼───────────────────────┼───────────────────────┘
                                 │                       │
                    ┌─────────────────┐    ┌─────────────────┐
                    │   PostgreSQL    │    │   Apache Kafka  │
                    │                 │    │                 │
                    │ • Пользователи  │    │ • order-events  │
                    │ • Товары        │    │ • Event-driven  │
                    │ • Заказы        │    │ • Асинхронность │
                    └─────────────────┘    └─────────────────┘
```

## 🚀 Основные возможности

### 👤 User Service (Порт: 8081)

- **Регистрация пользователей** с валидацией
- **Авторизация** через JWT токены
- **Управление профилями** (поставщики/клиенты)
- **Валидация токенов** для других сервисов

### 📦 Product Service (Порт: 8082)

- **Каталог товаров** с полным CRUD
- **Поиск и фильтрация** товаров
- **Управление товарами** поставщиками
- **Информация о поставщиках**

### 🛍️ Order Service (Порт: 8083)

- **Создание заказов** клиентами
- **Управление заказами** поставщиками
- **Отслеживание статусов** заказов
- **История заказов** для пользователей
- **Event-driven архитектура** через Kafka

### 📧 Notification Service (Порт: 8084)

- **Email уведомления** через SMTP
- **SMS уведомления** (мок)
- **Push уведомления** (мок)
- **Асинхронная обработка** событий заказов
- **Персонализированные сообщения**

## 🛠️ Технологический стек

### Backend

- **Go 1.21+** - основной язык программирования
- **Gorilla Mux** - HTTP роутинг
- **PostgreSQL** - основная база данных
- **Apache Kafka** - message broker для событий
- **JWT** - аутентификация и авторизация

### Инфраструктура

- **Docker Compose** - контейнеризация и оркестрация
- **Kafka UI** - веб-интерфейс для управления Kafka
- **Zookeeper** - координация Kafka кластера

### Библиотеки

- **github.com/IBM/sarama** - Kafka клиент для Go
- **github.com/golang-jwt/jwt** - JWT токены
- **gopkg.in/gomail.v2** - отправка email
- **github.com/lib/pq** - PostgreSQL драйвер

## 📁 Структура проекта

```
marketplace/
├── User_Service/                 # Сервис пользователей
│   ├── cmd/main.go              # Точка входа
│   ├── internal/
│   │   ├── api/                 # HTTP handlers
│   │   ├── models/              # Модели данных
│   │   └── repository/          # Работа с БД
│   └── go.mod
├── Product_Service/              # Сервис товаров
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── api/
│   │   ├── models/
│   │   └── repository/
│   └── go.mod
├── Order_Service/                # Сервис заказов
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── api/
│   │   ├── models/
│   │   ├── repository/
│   │   └── kafka/               # Kafka Producer
│   └── go.mod
├── Notification_Service/         # Сервис уведомлений
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── models/
│   │   ├── service/             # Email, SMS, Push сервисы
│   │   └── kafka/               # Kafka Consumer
│   └── go.mod
├── docker-compose.yml           # Docker конфигурация
└── README.md
```

## 🚀 Быстрый старт

### 1. Запуск инфраструктуры

```bash
# Запуск Kafka, Zookeeper и Kafka UI
docker-compose up -d

# Проверка статуса
docker-compose ps
```

### 2. Создание топика Kafka

Откройте [Kafka UI](http://localhost:8080) и создайте топик:

- **Название**: `order-events`
- **Партиции**: 3
- **Репликации**: 1

### 3. Запуск сервисов

```bash
# Терминал 1 - User Service
cd User_Service && go run cmd/main.go

# Терминал 2 - Product Service
cd Product_Service && go run cmd/main.go

# Терминал 3 - Order Service
cd Order_Service && go run cmd/main.go

# Терминал 4 - Notification Service
cd Notification_Service && go run cmd/main.go
```

## 🧪 Тестирование API

### Регистрация пользователей

```bash
# Регистрация поставщика
curl -X POST http://localhost:8081/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"supplier1","email":"supplier1@example.com","password":"password123","role":"supplier"}'

# Регистрация клиента
curl -X POST http://localhost:8081/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"client1","email":"client1@example.com","password":"password123","role":"client"}'
```

### Авторизация

```bash
# Получение JWT токена
curl -X POST http://localhost:8081/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"supplier1@example.com","password":"password123"}'
```

### Создание товара

```bash
# Создание товара (используйте токен из авторизации)
curl -X POST http://localhost:8082/api/product/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name":"iPhone 15","description":"Latest iPhone","price":999.99,"supplier_id":1}'
```

### Создание заказа

```bash
# Создание заказа (используйте токен клиента)
curl -X POST http://localhost:8083/api/order/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer CLIENT_JWT_TOKEN" \
  -d '{"product_name":"iPhone 15","product_id":1,"supplier_id":1}'
```

## 📊 Event-Driven Architecture

### Поток событий

```
1. Клиент создает заказ → Order Service
2. Order Service публикует событие "order_created" → Kafka
3. Notification Service получает событие → Kafka
4. Notification Service отправляет уведомления:
   - Email поставщику о новом заказе
   - Email клиенту о подтверждении
   - SMS и Push уведомления
```

### Типы событий

- **`order_created`** - новый заказ создан
- **`order_status_updated`** - статус заказа изменен

## 🔧 Конфигурация

### Порты сервисов

- **User Service**: 8081
- **Product Service**: 8082
- **Order Service**: 8083
- **Notification Service**: 8084
- **Kafka UI**: 8080

### База данных

- **PostgreSQL**: localhost:5432
- **База**: marketplace
- **Пользователь**: postgres
- **Пароль**: password

### Kafka

- **Broker**: localhost:9092
- **Zookeeper**: localhost:2181
- **Топик**: order-events

## 📧 Настройка уведомлений

### Email (Gmail)

```go
// В Notification_Service/cmd/main.go
emailService := service.NewEmailService(
    "smtp.gmail.com",           // SMTP host
    587,                        // SMTP port
    "your-email@gmail.com",     // SMTP username
    "your-app-password",        // SMTP password (App Password)
    "noreply@marketplace.com",  // From email
)
```

### SMS и Push (мок)

По умолчанию используются мок-сервисы для демонстрации. Для продакшена замените на реальные API:

- **SMS**: Twilio, Nexmo, SMS.ru
- **Push**: Firebase Cloud Messaging, Apple Push Notification Service

## 🐳 Docker

### Запуск всей инфраструктуры

```bash
docker-compose up -d
```

### Остановка

```bash
docker-compose down
```

### Логи

```bash
docker-compose logs -f kafka
```

## 📈 Мониторинг

### Kafka UI

- **URL**: http://localhost:8080
- **Функции**: просмотр топиков, сообщений, consumer groups

### Логи сервисов

Каждый сервис логирует:

- HTTP запросы
- Ошибки базы данных
- Kafka события
- Уведомления

## 🔒 Безопасность

- **JWT токены** для аутентификации
- **TLS шифрование** для SMTP
- **Валидация входных данных**
- **Обработка ошибок**

## 🚀 Развертывание

### Локальная разработка

1. Установите Go 1.21+
2. Установите Docker и Docker Compose
3. Клонируйте репозиторий
4. Запустите `docker-compose up -d`
5. Запустите сервисы по очереди

### Продакшен

- Используйте Kubernetes или Docker Swarm
- Настройте мониторинг (Prometheus, Grafana)
- Настройте логирование (ELK Stack)
- Используйте внешние базы данных и Kafka кластеры

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📝 Лицензия

MIT License

## 👥 Авторы

- **Разработчик**: [Ваше имя]
- **Архитектура**: Микросервисы + Event-Driven
- **Технологии**: Go, PostgreSQL, Kafka, Docker

---

**Marketplace** - современная платформа электронной коммерции с микросервисной архитектурой! 🛒✨
