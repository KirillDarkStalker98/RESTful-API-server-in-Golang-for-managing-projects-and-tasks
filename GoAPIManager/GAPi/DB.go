package GoAPIManager

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	// Инициализация базы данных
	env := "GAPi/DataBase.env"
	errenv := godotenv.Load(env) //".env"
	if errenv != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	//Для Windows
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	//Для докера
	/*connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	*/
	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных")
	}

	fmt.Println("База данных успешно подключена!")

	// Автоматическая миграция
	db.AutoMigrate(&User{}, &Project{}, &Task{})
	fmt.Println("Миграция базы данных выполнена успешно!")
}
