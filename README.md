# Room Booking Service

Сервис бронирования переговорок на Go.

## Технологии

- Go 1.21, Gorilla Mux, GORM, PostgreSQL 15
- JWT, Swagger (swaggo/swag), Docker Compose

## Быстрый запуск

```bash
docker-compose up --build


Метод	Эндпоинт	Описание	Роль
POST	/dummyLogin	Тестовый JWT	–
POST	/register	Регистрация	–
POST	/login	Вход	–
GET	/rooms/list	Список переговорок	admin/user
POST	/rooms/create	Создать переговорку	admin
POST	/rooms/{id}/schedule/create	Создать расписание	admin
GET	/rooms/{id}/slots/list?date=YYYY-MM-DD	Свободные слоты	admin/user
POST	/bookings/create	Создать бронь	user
GET	/bookings/my	Мои брони	user
POST	/bookings/{id}/cancel	Отменить бронь	user
GET	/bookings/list	Все брони (пагинация)	admin
GET	/_info	Health check	–

Генерация слотов on‑the‑fly – слоты создаются при первом запросе на конкретную дату. Это экономит БД и соответствует сценарию (99.9% запросов на ближайшие 7 дней).

Расписание неизменяемо – создаётся один раз, изменение невозможно (требование задания).

Идемпотентная отмена – повторная отмена брони возвращает 200 OK.

JWT содержит user_id и role. При создании брони user_id извлекается из токена.

Conference link – опциональный параметр createConferenceLink генерирует мок‑ссылку.

