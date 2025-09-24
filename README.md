# E-Commerce Platform

Полнофункциональная e-commerce платформа с современным frontend на Next.js и высокопроизводительным backend на Go.

## 🚀 Быстрый старт

### Одна команда для запуска всего:

```bash
make setup && make up
```

## 📋 Что включено

- **Frontend**: Next.js 14 с современным UI/UX
- **Backend**: Go + Gin с PostgreSQL
- **База данных**: PostgreSQL 15
- **Кэш**: Redis 7
- **Прокси**: Nginx
- **Контейнеризация**: Docker + Docker Compose

## 🛠️ Установка

### 1. Клонирование и настройка

```bash
git clone <repository-url>
cd E-Commerce
make setup
```

### 2. Настройка переменных окружения

Отредактируйте файл `.env`:

```bash
# Обязательные настройки
DB_PASSWORD=your-super-secure-password
JWT_SECRET=your-super-secret-jwt-key

# Опциональные настройки
BACKEND_PORT=5000
FRONTEND_PORT=3000
HTTP_PORT=80
```

### 3. Запуск

```bash
# Запуск всех сервисов
make up

# Или пошагово
make build
docker-compose up -d
```

## 🌐 Доступ к приложению

- **Frontend**: http://localhost
- **Backend API**: http://localhost/api
- **Health Check**: http://localhost/health

## 📊 Управление

### Основные команды

```bash
make up          # Запуск всех сервисов
make down        # Остановка всех сервисов
make restart     # Перезапуск
make logs        # Просмотр логов
make status      # Статус сервисов
make clean       # Очистка (удаляет данные!)
```

### Разработка

```bash
make dev-backend   # Backend в режиме разработки
make dev-frontend  # Frontend в режиме разработки
```

### База данных

```bash
make migrate-up    # Применить миграции
make migrate-down  # Откатить миграции
make db-shell      # Подключиться к БД
make db-backup     # Создать бэкап
```

## 🔒 Безопасность

### Что защищено:

- ✅ Все пароли в переменных окружения
- ✅ JWT секреты не в коде
- ✅ SSL/TLS для продакшена
- ✅ Security headers
- ✅ CORS настройки
- ✅ .env файлы в .gitignore

### Для продакшена:

1. **Измените все пароли** в `.env`
2. **Используйте HTTPS** (SSL сертификаты)
3. **Настройте firewall**
4. **Регулярно обновляйте** зависимости

## 🏗️ Архитектура

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Nginx     │────│  Frontend   │────│   Backend   │
│  (Port 80)  │    │ (Next.js)   │    │    (Go)     │
└─────────────┘    └─────────────┘    └─────────────┘
                           │                   │
                           │                   │
                    ┌─────────────┐    ┌─────────────┐
                    │   Redis     │    │ PostgreSQL  │
                    │   (Cache)   │    │ (Database)  │
                    └─────────────┘    └─────────────┘
```

## 📁 Структура проекта

```
E-Commerce/
├── frontend/           # Next.js приложение
├── backend-go/         # Go API сервер
├── nginx/             # Nginx конфигурация
├── scripts/           # Утилиты и скрипты
├── docker-compose.yml # Основная конфигурация
├── docker-compose.prod.yml # Продакшен конфигурация
├── Makefile          # Команды управления
└── README.md         # Документация
```

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_PASSWORD` | Пароль БД | **Обязательно** |
| `JWT_SECRET` | JWT секрет | **Обязательно** |
| `BACKEND_PORT` | Порт backend | `5000` |
| `FRONTEND_PORT` | Порт frontend | `3000` |
| `HTTP_PORT` | Порт nginx | `80` |

### Порты

- **80**: Nginx (основной доступ)
- **3000**: Frontend (прямой доступ)
- **5000**: Backend API (прямой доступ)
- **5432**: PostgreSQL
- **6379**: Redis

## 🚀 Продакшен

### 1. Подготовка

```bash
# Настройте продакшен переменные
cp env.example .env
# Отредактируйте .env с продакшен настройками

# Сгенерируйте SSL сертификаты
chmod +x scripts/generate-ssl.sh
./scripts/generate-ssl.sh
```

### 2. Запуск

```bash
make prod-build
make prod-up
```

### 3. Мониторинг

```bash
make status    # Статус сервисов
make logs      # Логи
make health    # Проверка здоровья
```

## 🐛 Отладка

### Просмотр логов

```bash
# Все сервисы
make logs

# Конкретный сервис
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Подключение к контейнерам

```bash
# Backend
docker-compose exec backend sh

# Frontend
docker-compose exec frontend sh

# База данных
make db-shell
```

### Очистка

```bash
make clean  # Удаляет ВСЕ данные!
```

## 📈 Производительность

- **Go backend**: Высокая производительность и низкое потребление памяти
- **Next.js frontend**: SSR и оптимизация
- **PostgreSQL**: Индексированная БД
- **Redis**: Кэширование
- **Nginx**: Обратный прокси и сжатие

## 🤝 Разработка

### Добавление новых функций

1. **Backend**: Добавьте handlers в `backend-go/internal/handlers/`
2. **Frontend**: Добавьте компоненты в `frontend/components/`
3. **База данных**: Создайте миграции в `backend-go/migrations/`

### Тестирование

```bash
# Backend тесты
cd backend-go && make test

# Frontend тесты
cd frontend && npm test
```

## 📄 Лицензия

MIT License

## 🆘 Поддержка

- Создайте issue в репозитории
- Проверьте логи: `make logs`
- Проверьте статус: `make status`

---

**Готово к использованию!** 🎉