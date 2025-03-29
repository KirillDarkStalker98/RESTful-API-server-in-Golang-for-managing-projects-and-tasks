package main

import (
	GoAPIManager "GoAPIManager/GAPi"
	_ "GoAPIManager/docs"
)

// @title           RESTful API-сервер на Golang
// @version         1.0
// @description     API для управления проектами и задачами, позволяющее пользователям создавать проекты, управлять задачами внутри них, а также изменять их свойства.
// @host            localhost:8080
// @BasePath        /

func main() {
	GoAPIManager.Controller()
}
