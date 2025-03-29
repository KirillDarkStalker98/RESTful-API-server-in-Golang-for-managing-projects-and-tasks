// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "description": "Аутентификация пользователя по имени и паролю и выдача access-токена",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Аутентификация"
                ],
                "summary": "Аутентификация пользователя",
                "parameters": [
                    {
                        "description": "Данные пользователя (имя и пароль)",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/GoAPIManager.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная аутентификация. Возвращает access-токен",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректные входные данные",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неверное имя пользователя или пароль",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка при генерации access-токена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/projects": {
            "post": {
                "description": "Создаёт новый проект, привязывая его к пользователю, авторизованному через JWT-токен",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Создание проекта",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Данные проекта",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/GoAPIManager.Project"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Проект успешно создан",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Некорректный ввод данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Необходим авторизационный токен или неверный формат токена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка базы данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/projects/:id/tasks/:task_id": {
            "put": {
                "description": "Обновляет задачу в проекте по ID, с проверкой обязательных полей и значений",
                "tags": [
                    "Задачи"
                ],
                "summary": "Обновление задачи",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID задачи",
                        "name": "task_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Обновленные данные задачи",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/GoAPIManager.Task"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Информация об обновленной задаче",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Ошибка валидации данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Задача не найдена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/projects/{id}": {
            "get": {
                "description": "Возвращает данные проекта по его ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Получение проекта",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Проект успешно найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Некорректный ID проекта",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Проект не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка базы данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Обновляет данные существующего проекта по его ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Обновление проекта",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Данные для обновления проекта",
                        "name": "project",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/GoAPIManager.Project"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Проект успешно обновлён",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Неверный запрос или некорректные данные",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Проект не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка базы данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаляет проект по указанному ID",
                "tags": [
                    "Проекты"
                ],
                "summary": "Удаление проекта",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Проект удален",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Проект не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/projects/{id}/download": {
            "get": {
                "description": "Скачивает файл, связанный с указанным проектом, если он существует на сервере",
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Скачивание файла проекта",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Файл проекта",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Файл не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/projects/{id}/tasks": {
            "get": {
                "description": "Получает список задач проекта с возможностью фильтрации по статусу, дедлайну и приоритету",
                "tags": [
                    "Задачи"
                ],
                "summary": "Получение задач проекта",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Статус задачи (In_Progress, Done, In_Line)",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Дедлайн задачи (формат: YYYY-MM-DD)",
                        "name": "deadline",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Приоритет задачи (High, Medium, Low)",
                        "name": "priority",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список задач",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Ошибка валидации данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Задачи не найдены",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Создает новую задачу для проекта с привязкой к пользователю",
                "tags": [
                    "Задачи"
                ],
                "summary": "Создание задачи",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Данные задачи",
                        "name": "task",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/GoAPIManager.Task"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Задача успешно создана",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Ошибка валидации данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/projects/{id}/upload": {
            "post": {
                "description": "Загружает файл для указанного проекта и сохраняет путь к файлу в базе данных",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Загрузка файла к проекту",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID проекта",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Файл для загрузки (до 100MB)",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Файл успешно загружен",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос или ошибка загрузки файла",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Проект не найден",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка при сохранении файла или обновлении базы данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/refresh/:id": {
            "post": {
                "description": "Проверяет refresh token, генерирует новый access и refresh токены и сохраняет новый refresh token в базе данных.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Аутентификация"
                ],
                "summary": "Обновление access и refresh токенов",
                "parameters": [
                    {
                        "description": "Тело запроса с refresh токеном",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/GoAPIManager.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Новые access и refresh токены",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректные входные данные",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Неверный refresh token",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "Регистрирует нового пользователя в системе.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Аутентификация"
                ],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "description": "Данные для регистрации",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/GoAPIManager.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пользователь успешно зарегистрирован",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Некорректные данные или пользователь уже существует",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/tasks/{task_id}": {
            "delete": {
                "description": "Удаляет задачу из базы данных по ID",
                "tags": [
                    "Задачи"
                ],
                "summary": "Удаление задачи",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID задачи",
                        "name": "task_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Сообщение об успешном удалении задачи",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Ошибка, если ID задачи некорректен",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Задача не найдена",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/user/projects": {
            "get": {
                "description": "Возвращает список проектов, созданных текущим пользователем",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Проекты"
                ],
                "summary": "Получение проектов пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer токен",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Проекты успешно найдены",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "401": {
                        "description": "Неавторизованный доступ или неверный токен",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "У пользователя нет проектов",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Ошибка базы данных",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "GoAPIManager.Project": {
            "type": "object",
            "properties": {
                "assignee_id": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "GoAPIManager.Task": {
            "type": "object",
            "required": [
                "priority",
                "status",
                "title"
            ],
            "properties": {
                "assignee_id": {
                    "type": "integer"
                },
                "deadline": {
                    "type": "string"
                },
                "description": {
                    "type": "string",
                    "maxLength": 500
                },
                "id": {
                    "type": "integer"
                },
                "priority": {
                    "type": "string",
                    "enum": [
                        "High",
                        "Medium",
                        "Low"
                    ]
                },
                "projectID": {
                    "type": "integer"
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "In_Progress",
                        "Done",
                        "In_Line"
                    ]
                },
                "title": {
                    "type": "string",
                    "maxLength": 255
                }
            }
        },
        "GoAPIManager.User": {
            "type": "object",
            "required": [
                "role"
            ],
            "properties": {
                "id": {
                    "type": "integer"
                },
                "password": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                },
                "role": {
                    "type": "string",
                    "enum": [
                        "User",
                        "Admin"
                    ]
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "RESTful API-сервер на Golang",
	Description:      "API для управления проектами и задачами, позволяющее пользователям создавать проекты, управлять задачами внутри них, а также изменять их свойства.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
