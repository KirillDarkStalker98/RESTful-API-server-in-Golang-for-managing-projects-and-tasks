# RESTful-API-server-in-Golang-for-managing-projects-and-tasks (RESTful API-сервер на Golang для управления проектами и задачами)

## Функционал API

* Регистация пользователя (Выдаёт RefreshToken, чтобы всегда можно было получить новый AccessToken)

* Вход (Выдаёт AccessToken для выполнения операций)

* Рефреш токена (Если устареет)

* Создание, удаление и обновление проекта

* Получение всех проектов конкретного пользователя или получение конкретного проекта

* Загрузка файлов в проект

* Создание, удаление, обновление и получение задач к проекту

Так же добавлен эндпоинт `/docs` для просмотра документации. 

## Что нужно для запуска на windows:

1. PostgreSQL
							
2. Язык Golang

* В начале для работы с сервисом в PostgreSQL необходимо создать базу данных с названием "gapim", далее миграции в коде сами справятся с созданием таблиц и связей наша база готова к работе с сервисом (Пользователь postgres и пароль в .env файле).

* Теперь когда база данных создана можно запустить сервер, переходим по расположению файла main.go в консоли (куда вы его скачали) и запускаем (Команда для запуска: "go run main.go").
  
* После запуска нужно дождаться сообщений (База данных успешно подключена! и Миграция базы данных выполнена успешно!) и можно работать👍.

## Что нужно для запуска через Docker:

Для работы с Docker всё немного проще. Сначала убедитесь, что Docker установлен на вашем компьютере. Затем выполните следующие шаги:

* В файле DB.go нужно поменять код "Для Windows", на код "Для докера" который закомментирован.

* Перейдите в директорию, куда вы скачали файлы проекта используя консоль.

* В командной строке выполните команду: ("docker-compose build"), эта команда соберет образы, указанные в файле docker-compose.yml.

* После успешной сборки образов вы можете запустить контейнеры с помощью команды: ("docker-compose up") либо с помощью интерфеса через Docker Desktop, команда чтобы выключить ("docker-compose down").

## Curl запросы:

### 1. Регистация пользователя

* Запрос: curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d "{ \"username\": \"User1\", \"password\": \"wordPass243\", \"role\": \"admin\" }"
   
* Ответ: {"message":"Пользователь успешно зарегистрирован","Ваш RefreshToken, сохраните его для того чтобы его можно было обменять на новый AccessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJ1c2VybmFtZSI6IlVzZXIxIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzQzNzY3MzI1fQ.4XDMUyDntCj7IMwXEuKWS4IIV9SO71FrVWDZ2wtTwqQ"}

### 2. Вход

* Запрос: curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{ \"username\": \"User1\", \"password\": \"wordPass243\" }"
   
* Ответ: {"AccessToken для всех последующих операций":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzIwNTc5Nn0._9dfIwi96mFNuMBb8W6L4UufXf6CJK5TtLFnVgbqkB8","message":"Вы успешно вошли"}

### 3. Создание проекта

* Запрос: curl -X POST http://localhost:8080/projects -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzIwNTc5Nn0._9dfIwi96mFNuMBb8W6L4UufXf6CJK5TtLFnVgbqkB8" -H "Content-Type: application/json" -d "{ \"name\": \"Project A\", \"description\": \"Test project\" }"
   
* Ответ: {"Projects:":{"ID":19,"Name":"Project A","Description":"Test project","CreatedAt":"2025-03-28T16:52:22.5505805+05:00","assignee_id":17,"message":"Проект успешно создан"}


