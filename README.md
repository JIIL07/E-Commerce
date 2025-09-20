# 🛒 E-Commerce Platform

Полнофункциональная платформа электронной коммерции с современным дизайном, корзиной покупок, системой платежей через Stripe, админ-панелью и мобильной адаптацией.

## 🚀 Особенности

- **Современный дизайн** - Красивый и отзывчивый интерфейс с Tailwind CSS
- **Корзина покупок** - Полнофункциональная корзина с управлением состоянием
- **Платежи Stripe** - Безопасная обработка платежей
- **Админ-панель** - Управление продуктами, заказами и пользователями
- **Мобильная адаптация** - Полная поддержка мобильных устройств
- **Аутентификация** - Безопасная система входа и регистрации
- **База данных** - PostgreSQL с Prisma ORM
- **TypeScript** - Полная типизация для надежности

## 📁 Структура проекта

```
ecommerce-platform/
├── frontend/          # Next.js фронтенд приложение
├── backend/           # Express.js API сервер
├── admin/             # Next.js админ-панель
├── database/          # Prisma схема и миграции
├── package.json       # Корневые скрипты
├── .env.example       # Пример переменных окружения
└── README.md          # Документация
```

## 🛠 Технологии

### Frontend
- **Next.js 14** - React фреймворк
- **TypeScript** - Типизация
- **Tailwind CSS** - Стилизация
- **Zustand** - Управление состоянием
- **React Hook Form** - Формы
- **Framer Motion** - Анимации
- **Stripe Elements** - Платежи

### Backend
- **Express.js** - API сервер
- **TypeScript** - Типизация
- **Prisma** - ORM для базы данных
- **JWT** - Аутентификация
- **Stripe** - Обработка платежей
- **Multer** - Загрузка файлов
- **Cloudinary** - Хранение изображений

### Database
- **PostgreSQL** - Основная база данных
- **Prisma** - ORM и миграции

### Admin Panel
- **Next.js 14** - Админ-панель
- **Recharts** - Графики и аналитика
- **React Table** - Таблицы данных

## 🚀 Быстрый старт

### Предварительные требования

- Node.js 18+ 
- PostgreSQL 14+
- npm или yarn

### Установка

1. **Клонируйте репозиторий**
   ```bash
   git clone <repository-url>
   cd ecommerce-platform
   ```

2. **Установите зависимости**
   ```bash
   npm run install:all
   ```

3. **Настройте переменные окружения**
   ```bash
   cp .env.example .env
   ```
   
   Отредактируйте `.env` файл с вашими настройками:
   - Настройте подключение к базе данных
   - Добавьте ключи Stripe
   - Настройте email для уведомлений
   - Добавьте ключи Cloudinary

4. **Настройте базу данных**
   ```bash
   npm run db:generate
   npm run db:migrate
   npm run db:seed
   ```

5. **Запустите проект**
   ```bash
   npm run dev
   ```

### Доступные URL

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:5000
- **Admin Panel**: http://localhost:3001
- **Database Studio**: http://localhost:5555 (после `npm run db:studio`)

## 📋 Доступные скрипты

### Корневые скрипты

```bash
# Разработка (запуск всех сервисов)
npm run dev

# Сборка всех проектов
npm run build

# Запуск в продакшене
npm run start

# Установка всех зависимостей
npm run install:all

# Линтинг всех проектов
npm run lint

# Проверка типов
npm run type-check

# Работа с базой данных
npm run db:generate    # Генерация Prisma клиента
npm run db:migrate     # Применение миграций
npm run db:studio      # Открытие Prisma Studio
npm run db:seed        # Заполнение тестовыми данными
```

### Индивидуальные скрипты

```bash
# Frontend
npm run dev:frontend
npm run build:frontend
npm run start:frontend

# Backend
npm run dev:backend
npm run build:backend
npm run start:backend

# Admin
npm run dev:admin
npm run build:admin
npm run start:admin
```

## 🗄 База данных

### Схема

- **Users** - Пользователи системы
- **Products** - Товары
- **Categories** - Категории товаров
- **Orders** - Заказы
- **OrderItems** - Элементы заказов
- **Cart** - Корзина покупок
- **Reviews** - Отзывы о товарах

### Миграции

```bash
# Создание новой миграции
cd database
npx prisma migrate dev --name migration_name

# Применение миграций в продакшене
npm run db:deploy
```

## 🔐 Аутентификация

Система использует NextAuth.js для аутентификации с поддержкой:
- Email/пароль
- OAuth провайдеров (Google, GitHub)
- JWT токенов

## 💳 Платежи

Интеграция со Stripe включает:
- Обработка платежей
- Webhooks для подтверждения
- Поддержка различных валют
- Безопасное хранение данных карт

## 📱 Мобильная адаптация

- Responsive дизайн для всех устройств
- Touch-friendly интерфейс
- Оптимизация для мобильных браузеров
- PWA поддержка

## 🧪 Тестирование

```bash
# Запуск тестов
npm test

# Тесты с покрытием
npm run test:coverage
```

## 🚀 Деплой

### Vercel (Frontend/Admin)
```bash
# Установите Vercel CLI
npm i -g vercel

# Деплой
vercel --prod
```

### Railway/Heroku (Backend)
```bash
# Настройте переменные окружения
# Деплой через Git
git push heroku main
```

## 📝 API Документация

API документация доступна по адресу: `http://localhost:5000/api-docs`

### Основные эндпоинты

- `GET /api/products` - Список товаров
- `POST /api/orders` - Создание заказа
- `GET /api/orders/:id` - Получение заказа
- `POST /api/payments` - Обработка платежа

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для новой функции
3. Внесите изменения
4. Добавьте тесты
5. Создайте Pull Request

## 📄 Лицензия

MIT License

## 🆘 Поддержка

Если у вас есть вопросы или проблемы:
- Создайте Issue в GitHub
- Обратитесь к документации
- Проверьте FAQ

---

**Создано с ❤️ для современной электронной коммерции**
