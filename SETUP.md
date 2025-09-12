# 🚀 Инструкция по установке и настройке

## Быстрый старт

### 1. Установка зависимостей

```bash
# Установка Go (если не установлен)
# Скачайте с https://golang.org/dl/

# Установка Node.js (если не установлен)
# Скачайте с https://nodejs.org/

# Установка Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. Клонирование и установка

```bash
git clone <repository-url>
cd todo-wails-go

# Установка Go зависимостей
go mod tidy

# Установка frontend зависимостей
cd frontend
npm install
cd ..
```

### 3. Запуск приложения

```bash
# Режим разработки (рекомендуется для тестирования)
wails dev

# Или сборка и запуск
wails build
./build/bin/todo-wails-go
```

## 🗄️ Настройка PostgreSQL (опционально)

### Установка PostgreSQL

#### macOS (с Homebrew)
```bash
brew install postgresql
brew services start postgresql
```

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### Windows
Скачайте с https://www.postgresql.org/download/windows/

### Создание базы данных

```bash
# Подключение к PostgreSQL
psql -U postgres

# Создание базы данных
CREATE DATABASE todo_app;

# Создание пользователя (опционально)
CREATE USER todo_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE todo_app TO todo_user;

# Выход
\q
```

### Настройка переменных окружения

```bash
# macOS/Linux
export DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=todo_app sslmode=disable"

# Windows
set DATABASE_URL=host=localhost port=5432 user=postgres password=postgres dbname=todo_app sslmode=disable
```

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | Строка подключения к PostgreSQL | `host=localhost port=5432 user=postgres password=postgres dbname=todo_app sslmode=disable` |

### Настройки приложения

Приложение автоматически:
- Создает таблицы в PostgreSQL при первом запуске
- Переключается на in-memory хранилище если PostgreSQL недоступен
- Сохраняет настройки темы в localStorage

## 🐛 Решение проблем

### Ошибка подключения к PostgreSQL

```
Warning: Failed to connect to PostgreSQL: ...
Using in-memory storage instead
```

**Решение**: Убедитесь что PostgreSQL запущен и настройки подключения корректны.

### Ошибка сборки frontend

```
Error: npm install failed
```

**Решение**: 
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Ошибка Wails

```
wails: command not found
```

**Решение**: Убедитесь что Go установлен и `$GOPATH/bin` в `$PATH`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## 📱 Тестирование

### Функциональное тестирование

1. **Создание задач**
   - Введите заголовок задачи
   - Добавьте описание (опционально)
   - Выберите приоритет
   - Установите дату выполнения (опционально)
   - Нажмите "Add Task"

2. **Управление задачами**
   - Отметьте задачу как выполненную (чекбокс)
   - Отредактируйте задачу (кнопка "Edit")
   - Удалите задачу (кнопка "Delete")

3. **Фильтрация и сортировка**
   - Используйте фильтры по статусу и приоритету
   - Измените сортировку
   - Очистите фильтры

4. **Темы**
   - Переключите тему (кнопка в заголовке)
   - Перезапустите приложение - тема сохранится

### Тестирование производительности

- Создайте 100+ задач
- Протестируйте фильтрацию и сортировку
- Проверьте время отклика

## 🚀 Развертывание

### Сборка для разных платформ

```bash
# macOS
wails build -platform darwin/amd64

# Windows
wails build -platform windows/amd64

# Linux
wails build -platform linux/amd64
```

### Создание установщика

```bash
# macOS (требует подписи)
wails build -platform darwin/amd64 -package

# Windows (требует NSIS)
wails build -platform windows/amd64 -package
```

## 📊 Мониторинг

### Логи приложения

Логи выводятся в консоль при запуске в режиме разработки:
```bash
wails dev
```

### Мониторинг базы данных

```sql
-- Подключение к PostgreSQL
psql -U postgres -d todo_app

-- Просмотр таблиц
\dt

-- Просмотр задач
SELECT * FROM tasks;

-- Статистика
SELECT 
    status,
    priority,
    COUNT(*) as count
FROM tasks 
GROUP BY status, priority;
```

## 🔒 Безопасность

### Рекомендации для продакшена

1. **База данных**
   - Используйте сильные пароли
   - Настройте SSL/TLS
   - Ограничьте доступ по IP

2. **Приложение**
   - Не храните пароли в коде
   - Используйте переменные окружения
   - Регулярно обновляйте зависимости

3. **Данные**
   - Регулярно создавайте резервные копии
   - Шифруйте чувствительные данные
   - Логируйте действия пользователей

## 📞 Поддержка

При возникновении проблем:

1. Проверьте логи приложения
2. Убедитесь что все зависимости установлены
3. Проверьте настройки базы данных
4. Создайте issue в репозитории

---

**Примечание**: Приложение работает без PostgreSQL, используя in-memory хранилище, но для продакшена рекомендуется использовать PostgreSQL для надежности и производительности.
