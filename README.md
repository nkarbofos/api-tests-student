# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Структура проекта

Проект построен по слоям, тест вызывает Domain Layer, Domain Layer вызывает Client Layer:

- `tests/scenarios/` -- Test Layer, сами сценарии;
- `internal/managers/` -- Domain Layer: бизнес-действия, проверка статус-кода и разбор тела ответа;
- `internal/client/http/` -- Client Layer: HTTP-вызовы ручек;
- `internal/assertions/` -- общие проверки ответов;
- `internal/constants/` -- пути ручек и общие коды/тексты ошибок API;
- `internal/helpers/` -- api-runner, фикстуры, создание тестовых фотографий, polling;
- `internal/utils/` -- переменные окружения, логирование, генераторы.

## Покрытые сценарии

Часть 1. Добавление объявления:

- `tests/scenarios/advertisementCreate` -- позитивный сценарий: создание объявления со всеми параметрами, проверка через `GET /advertisement`, поиск по полному названию и получение фотографий через `GET /advertisements/{id}/photos`;
- `tests/scenarios/advertisementCreateNegative` -- негативные сценарии: пустой и слишком длинный title, отрицательная и слишком большая цена, отрицательное quantity, отсутствие фотографий и превышение их лимита, запрос без токена и с неверным токеном.

Часть 2. Работа с избранным:

- `tests/scenarios/favorites` -- полный цикл: добавление в избранное, появление в списке, удаление из избранного, повторное добавление, удаление объявления владельцем и проверка, что оно исчезло из избранного;
- `tests/scenarios/favoritesNegative` -- негативные сценарии: несуществующий `ad_id` при добавлении и удалении, повторное добавление, добавление собственного объявления, запросы без авторизации.

Тесты создают данные через API и удаляют их после себя. Для сценариев, где требуется объявление другого пользователя, регистрируется отдельный тестовый пользователь с уникальным email -- это исключает конфликты между прогонами.

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1 или выше
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
TEST_LOGIN=your_login
TEST_PASSWORD=your_password
```

8. Для запуска всех тестов выполните команду в корне проекта:

```bash
  go test -v ./...
```

Клиент автоматически повторяет запросы при временных ответах TestBoard `429/502/503/504`. Если сервис отвечает `503 Service Unavailable` из-за ограничения частоты запросов, запустите пакеты последовательно:

```bash
  go test -p 1 -v ./...
```

9. Отдельный пакет с тестами можно запустить командой:

```bash
  go test -v ./tests/scenarios/myAdvertisement/
```

## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```