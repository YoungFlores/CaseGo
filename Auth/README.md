# Краткая инструкция для запуска

## Быстрый старт

```bash
# 1. Настраиваем
cd Auth
cp .env.example .env
python keys/generate_rsa_keys.py

# 2. Запускаем сервисы
docker-compose up --build -d

# 3. Применяем миграции
make upgrade

# 4. Открываем документацию API
open http://localhost:8000/docs
```

Примеры работы с API предоставлены в документации.

## Команды управления

### Миграции
```bash
make create message="add_users_table"
make upgrade
make downgrade
```
