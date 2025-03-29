package GoAPIManager

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/time/rate"

	_ "GoAPIManager/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Лимиты запросов
var rateLimiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	// Если для IP еще нет лимитера, создаем новый
	if _, exists := rateLimiters[ip]; !exists {
		// 100 запросов в минуту
		limiter := rate.NewLimiter(rate.Every(60*time.Second), 100)
		rateLimiters[ip] = limiter
		// Удаляем лимитер через 10 минут (чтобы не хранить в памяти вечно)
		go func() {
			time.Sleep(10 * time.Minute)
			mu.Lock()
			delete(rateLimiters, ip)
			mu.Unlock()
		}()
	}
	return rateLimiters[ip]
}

// Middleware для ограничения запросов
func rateLimitMiddleware(c *gin.Context) {
	ip := c.ClientIP()
	limiter := getLimiter(ip)

	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		c.Abort()
		return
	}

	c.Next()
}

// Логирование запросов
func requestLoggerMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start)

	statusCode := c.Writer.Status()
	logRequest(c, statusCode, fmt.Sprintf("Processed in %v", duration))
	// Логирование ошибок базы данных, если есть
	if len(c.Errors) > 0 {
		for _, e := range c.Errors {
			if strings.Contains(e.Error(), "database") {
				log.Printf("[DB ERROR] %s", e.Error())
			}
		}
	}
}

// Логирование запросов и ошибок
func logRequest(c *gin.Context, statusCode int, message string) {
	// Получаем IP клиента
	clientIP := c.ClientIP()

	// Получаем токен (если есть)
	tokenString := c.GetHeader("Authorization")

	var username string

	// Если токен есть, ищем пользователя
	if tokenString != "" {
		const prefix = "Bearer "
		if len(tokenString) > len(prefix) && strings.HasPrefix(tokenString, prefix) {
			tokenString = tokenString[len(prefix):]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err == nil && token.Valid {
			username = claims.Username
		}
	}

	// Формируем строку лога

	if username != "" {
		log.Printf("[REQUEST] %s | %s | %s | %d | %s", username, c.Request.Method, c.Request.URL.Path, statusCode, message)
	} else {
		log.Printf("[REQUEST] %s | %s | %s | %d | %s", clientIP, c.Request.Method, c.Request.URL.Path, statusCode, message)
	}
}

func Controller() {
	initDB()
	// Открываем лог-файл (Мои логи)
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия log-файла: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	// Открываем или создаём файл логов (Дефолтные логи Gin)
	logFile2, err := os.OpenFile("serverGIN.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка открытия log-файла: %v", err)
	}

	// Перенаправляем стандартный логгер Gin в файл и в консоль
	gin.DefaultWriter = io.MultiWriter(logFile2, os.Stdout)

	// Инициализация роутера
	r := gin.Default()

	// Применяем Rate Limit Middleware ко всем маршрутам
	r.Use(requestLoggerMiddleware)
	r.Use(rateLimitMiddleware)
	// Настройка rate limit (например, 10 запросов в минуту на IP)

	// Маршруты для аутентификации
	r.POST("/register", registerUser)
	r.POST("/login", loginUser)
	r.POST("/refresh/:id", refreshToken)

	// Маршруты для проектов
	auth := r.Group("/")
	auth.Use(authMiddleware)
	auth.POST("/projects", createProject)
	auth.GET("/user/projects", getUserProjects)
	auth.GET("/projects/:id", getProject)
	auth.PUT("/projects/:id", updateProject)
	auth.DELETE("/projects/:id", deleteProject)
	auth.POST("/projects/:id/upload", uploadProjectFile)
	auth.GET("/projects/:id/download", downloadProjectFile)

	// Маршруты для задач
	auth.POST("/projects/:id/tasks", createTask)
	auth.GET("/projects/:id/tasks", getTasks)
	auth.PUT("/projects/:id/tasks/:task_id", updateTask)
	auth.DELETE("/projects/:id/tasks/:task_id", deleteTask)

	// Эндпоинт для документации
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	r.Run(":8080")

	log.Println("Сервер работает на порту 8080")
	// Запуск сервера с логированием ошибок
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
