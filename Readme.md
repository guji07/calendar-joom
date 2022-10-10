# Calender-joom

Сервис календарь, позволяющий:
- создавать пользователей
- создавать встречи
- Получать детали встречи
- приглашать пользователей во встречи
- принимать и отклонять приглашения
- получать список встреч для пользователя
- получать ближайшее окно для встречи определенной продолжительности учитывая календари заданных пользователей
- все встречи поддерживают повторы по формату RRule iCalendar RFC


## _Запуск_

Следующая команда поднимает базу данных(postgresql) в докер-контейнере и запускает приложение

Миграция запускается самим приложением
```bash
make
```

Запуск линтера:
```bash
make lint
```

Запуск тестов:
```bash
make test
```
# Описание http api

##CreateUser
```bash
curl --location --request POST 'localhost:8080/users' \
--header 'Content-Type: application/json' \
--data-raw '{
"login" : "login3",
"first_name": "name",
"last_name": "lastname"
}'
```
##GetUserEvents
```bash
curl --location --request GET 'localhost:8080/users/3/events/?from=2021-10-09T20:43:00Z&to=2023-10-09T20:43:00Z'
```
##CreateEvent
```bash
curl --location --request POST 'localhost:8080/events' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "event_name",
    "author": 3,
    "repeatable": true,
    "repeat_options": "DTSTART;TZID=Europe/Moscow:20221009T220000\nFREQ=HOURLY;INTERVAL=1;COUNT=10",
    "begin_time": "2022-10-09T20:43:00+00:00",
    "end_time":   "2022-10-09T20:44:00+00:00",
    "is_private": false,
    "details": "details of the meeting",
    "invited_users": [1]
}'
```
##GetEvent
```bash
curl --location --request GET 'localhost:8080/events/1?user_id=2'
```
##RespondOnEvent
```bash
curl --location --request GET 'localhost:8080/events/5/respond?user_id=3&accept=true'
```
##FindWindowForEvent
```bash
curl --location --request POST 'localhost:8080/events/window_by_users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "users_ids": [1],
    "duration": 3
}'
```
