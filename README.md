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

## Примеры запросы:

### 1. Регистация пользователя

* Запрос: curl -X POST http://localhost:8080/register -H "Content-Type: application/json" -d "{ \"username\": \"User1\", \"password\": \"wordPass243\", \"role\": \"admin\" }"
   
* Ответ: {"message":"Пользователь успешно зарегистрирован","Ваш RefreshToken, сохраните его для того чтобы его можно было обменять на новый AccessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJ1c2VybmFtZSI6IlVzZXIxIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzQzNzY3MzI1fQ.4XDMUyDntCj7IMwXEuKWS4IIV9SO71FrVWDZ2wtTwqQ"}

### 2. Вход

* Запрос: curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{ \"username\": \"User1\", \"password\": \"wordPass243\" }"
   
* Ответ: {"AccessToken для всех последующих операций":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzIwNTc5Nn0._9dfIwi96mFNuMBb8W6L4UufXf6CJK5TtLFnVgbqkB8","message":"Вы успешно вошли"}

### 3. Создание проекта

* Запрос: curl -X POST http://localhost:8080/projects -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzIwNTc5Nn0._9dfIwi96mFNuMBb8W6L4UufXf6CJK5TtLFnVgbqkB8" -H "Content-Type: application/json" -d "{ \"name\": \"Project A\", \"description\": \"Test project\" }"
   
* Ответ: {"Projects:":{"ID":19,"Name":"Project A","Description":"Test project","CreatedAt":"2025-03-28T16:52:22.5505805+05:00","assignee_id":17,"message":"Проект успешно создан"}

### 4. Получение всех проектов пользователя (по токену)

* Запрос: curl -X GET http://localhost:8080/user/projects -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Ответ: {"Projects:":[{"ID":19,"Name":"Project A","Description":"Test project","CreatedAt":"2025-03-28T16:52:22.55058+05:00","assignee_id":17}],"message":"Проекты успешно найдены"}

### 5. Получение конкретного проекта пользователя (по id проекта)

* Запрос: curl -X GET http://localhost:8080/projects/19 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Ответ: {"Projects:":{"ID":19,"Name":"Project A","Description":"Test project","CreatedAt":"2025-03-28T16:52:22.55058+05:00","assignee_id":17},"message":"Проект успешно найден"}

### 6. Обновление проекта

* Запрос: curl -X PUT http://localhost:8080/projects/19 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA" -H "Content-Type: application/json" -d "{ \"name\": \"Project B\", \"description\": \"Test project UPDATE\" }"

* Ответ: {"Projects:":{"ID":19,"Name":"Project B","Description":"Test project UPDATE","CreatedAt":"2025-03-28T16:52:22.55058+05:00","assignee_id":17},"message":"Проект успешно обновлён"}

### 7. Загрузка файлов в проект

* Запрос: curl -X POST http://localhost:8080/projects/19/upload -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA" -F "file=@O:\Test\fileproject.txt"

* Ответ: {"File_path":"uploads\\project_19_fileproject.txt","message":"Файл успешно загружен"}

### 8. Загрузка файлов из проекта

* Запрос: curl -X GET http://localhost:8080/projects/19/download -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Ответ: TEST SAVE 124 ТЕСТ СОХРАНЕНИЯ 123 ввв

### 9. Создание задачи

* Запрос: curl -X POST http://localhost:8080/projects/19/tasks -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA" -H "Content-Type: application/json" -d "{ \"title\": \"Task 1\", \"description\": \"Do something\", \"status\": \"In_Progress\", \"priority\": \"High\", \"deadline\": \"2025-04-12T00:00:00Z\" }"

* Ответ: {"Task":{"ID":16,"ProjectID":19,"title":"Task 1","description":"Do something","status":"In_Progress","priority":"High","deadline":"2025-04-12T00:00:00Z","assignee_id":17},"message":"Задача успешно создана"}

### 10. Получение всех задач для проекта

* Запрос: curl -X GET http://localhost:8080/projects/19/tasks -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Поиск по дате: curl -X GET http://localhost:8080/projects/19/tasks?deadline=2025-04-15 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Поиск по статусу: curl -X GET http://localhost:8080/projects/19/tasks?status=In_Progress -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Поиск по приоритету: curl -X GET http://localhost:8080/projects/19/tasks?priority=High -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Ответ: {"Задачи":[{"ID":16,"ProjectID":19,"title":"Task 1","description":"Do something","status":"In_Progress","priority":"High","deadline":"2025-04-12T05:00:00+05:00","assignee_id":17}]}

### 11. Обновление задачи

* Запрос: curl -X PUT http://localhost:8080/projects/19/tasks/16 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA" -H "Content-Type: application/json" -d "{ \"title\": \"Task UPDATE\", \"description\": \"Do something A\", \"status\": \"Done\", \"priority\": \"Low\" }"

* Ответ: {"Task":{"ID":16,"ProjectID":19,"title":"Task UPDATE","description":"Do something A","status":"Done","priority":"Low","deadline":"2025-04-12T05:00:00+05:00","assignee_id":17},"message":"Задача успешно обновлена"}

### 12. Удаление задачи

* Запрос: curl -X DELETE http://localhost:8080/projects/19/tasks/16 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Ответ: {"message":"Задача успешно удалена"}

### 13. Удаление проекта

* Запрос: curl -X DELETE http://localhost:8080/projects/19 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI3ODg3Mn0.luEQFrT4AQpFdXgDOJfmw5yVd2gQ9wOUWq-ZfROBGIA"

* Ответ: {"message":"Проект успешно удалён"}

### 14. Рефреш токена

* Запрос: curl -X POST http://localhost:8080/refresh/17 -H "Content-Type: application/json" -d "{ \"refreshtoken\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJ1c2VybmFtZSI6IlVzZXIxIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzQzNzY3MzI1fQ.4XDMUyDntCj7IMwXEuKWS4IIV9SO71FrVWDZ2wtTwqQ\" }"

* Ответ: {"AccessToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0MzI4MjkxNX0.0hklnN79jBlXAQOKIPtd-8yHYJk_YYB-vjDkyYm-_kc","RefreshToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNywidXNlcm5hbWUiOiJVc2VyMSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTc0Mzg0NDUxNX0.oTJU-0vRL64l9wkBai0Q3RFDUJeaIZQVxR6sBUC9RTo"}

### 14. Просмотр документации

* Запрос: Вставить в браузере (http://localhost:8080/docs/index.html#/)

* Ответ: ![image](https://github.com/user-attachments/assets/b2489ed4-6b05-48c1-bf7f-86b344c54261)

