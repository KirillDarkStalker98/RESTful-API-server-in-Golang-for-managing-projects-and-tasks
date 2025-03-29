package GoAPIManager

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const uploadDir = "uploads/"

var ctx = context.Background()

var validate = validator.New()

var jwtKey = []byte("secret_key")

// Models
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	Password     string `gorm:"not null"`
	Role         string `gorm:"not null" validate:"required,oneof=User Admin"`
	RefreshToken string `gorm:"column:refreshtoken"`
	//Добавить связь (Закомментировать после того как база данных создана, иначе будут при ответах вылазить ненужные строки)
	Tasks []Task `gorm:"foreignKey:AssigneeID"`
}

type Project struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	AssigneeID  uint `json:"assignee_id" gorm:"not null"`
	//Добавить связи (Закомментировать после того как база данных создана, иначе будут при ответах вылазить ненужные строки)
	Assignee User   `gorm:"foreignKey:AssigneeID"` // связь с User
	Tasks    []Task `gorm:"foreignKey:ProjectID"`  // связь с задачами
}

type Task struct {
	ID          uint      `gorm:"primaryKey"`
	ProjectID   uint      `gorm:"not null"`
	Title       string    `gorm:"not null" json:"title" validate:"required,max=255"`
	Description string    `json:"description" validate:"max=500"`
	Status      string    `gorm:"not null" json:"status" validate:"required,oneof=In_Progress Done In_Line"`
	Priority    string    `gorm:"not null" json:"priority" validate:"required,oneof=High Medium Low"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	AssigneeID  uint      `json:"assignee_id" gorm:"not null"`
	//Добавить связи (Закомментировать после того как база данных создана, иначе будут при ответах вылазить ненужные строки)
	Assignee User    `gorm:"foreignKey:AssigneeID;references:ID;constraint:OnDelete:SET NULL"` // связь с User
	Project  Project `gorm:"foreignKey:ProjectID;references:ID;constraint:OnDelete:CASCADE"`   // связь с Project
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		c.Abort()
		return
	}

	// Убираем "Bearer " из строки токена
	const prefix = "Bearer "
	if len(tokenString) > len(prefix) && strings.HasPrefix(tokenString, prefix) {
		tokenString = tokenString[len(prefix):]
	}

	tokenString = strings.TrimPrefix(tokenString, prefix)

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if claims.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		c.Abort()
		return
	}

	c.Set("id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("role", claims.Role)
	// Проверка, имеет ли пользователь доступ к управлению проектом/задачами
	// Карта маршрутов, где ключ — путь, а значение — список методов, требующих проверки
	protectedRoutes := map[string]map[string]bool{
		"/projects/:id": {
			http.MethodGet:    true,
			http.MethodPut:    true,
			http.MethodDelete: true,
		},
		"/projects/:id/upload": {
			http.MethodPost: true,
		},
		"/projects/:id/download": {
			http.MethodGet: true,
		},
		"/projects/:id/tasks": {
			http.MethodPost: true,
			http.MethodGet:  true,
		},
		"/projects/:id/tasks/:task_id": {
			http.MethodPut:    true,
			http.MethodDelete: true,
		},
	}
	// Проверяем, есть ли путь в списке защищенных
	if methods, exists := protectedRoutes[c.FullPath()]; exists {
		if methods[c.Request.Method] {
			projectID, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
				c.Abort()
				return
			}

			var project Project
			if err := db.First(&project, projectID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
				c.Abort()
				return
			}

			usery := claims.UserID

			if usery == 0 {
				// Обработка ошибки, если пользователь не найден
				fmt.Println("Ошибка:", usery)
			}

			userID := claims.UserID
			role := claims.Role

			if role != "admin" && userID != project.AssigneeID {
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}

		}
	}

	c.Next()
}

// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя в системе.
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param input body User true "Данные для регистрации"
// @Success 201 {object} map[string]interface{} "Пользователь успешно зарегистрирован"
// @Failure 400 {object} map[string]interface{} "Некорректные данные или пользователь уже существует"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /register [post]
func registerUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Проверка длины пароля
	if len(user.Password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
		return
	}

	// Валидация имени пользователя (пример: от 3 до 20 символов, только буквы и цифры)
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,20}$`, user.Username); !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username format"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Проверка, существует ли пользователь
	var existingUser User
	if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
			return
		}
	}

	// Генерация refresh-токена
	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Сохранение пользователя в БД
	user.Password = string(hashedPassword)
	user.RefreshToken = refreshToken

	//Транзакцией
	tx := db.Begin()
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	tx.Commit()

	// Отправка ответа
	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован", "Ваш RefreshToken, сохраните его для того чтобы его можно было обменять на новый AccessToken": refreshToken})
}

// @Summary Обновление access и refresh токенов
// @Description Проверяет refresh token, генерирует новый access и refresh токены и сохраняет новый refresh token в базе данных.
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body User true "Тело запроса с refresh токеном"
// @Success 200 {object} map[string]string "Новые access и refresh токены"
// @Failure 400 {object} map[string]string "Некорректные входные данные"
// @Failure 401 {object} map[string]string "Неверный refresh token"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /refresh/:id [post]
func refreshToken(c *gin.Context) {
	// Проверка, есть ли такой refreshToken в базе данных
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if user.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	if err := db.Where("refreshtoken = ?", user.RefreshToken).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Генерация нового accessToken для пользователя
	newAccessToken, err := generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	// Генерация нового refreshToken
	newRefreshToken, err := generateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	// Обновление refresh токена в базе данных с транзакцией
	tx := db.Begin()
	user.RefreshToken = newRefreshToken
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update refresh token"})
		return
	}
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"AccessToken": newAccessToken, "RefreshToken": newRefreshToken})
}

// Генерация refresh token
func generateRefreshToken(user User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // refreshToken срок жизни 7 дней
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// Функция для генерации access token
func generateAccessToken(user User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(), // accessToken срок жизни 12 часов
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// @Summary Аутентификация пользователя
// @Description Аутентификация пользователя по имени и паролю и выдача access-токена
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param input body User true "Данные пользователя (имя и пароль)"
// @Success 200 {object} map[string]string "Успешная аутентификация. Возвращает access-токен"
// @Failure 400 {object} map[string]string "Некорректные входные данные"
// @Failure 401 {object} map[string]string "Неверное имя пользователя или пароль"
// @Failure 500 {object} map[string]string "Ошибка при генерации access-токена"
// @Router /login [post]
func loginUser(c *gin.Context) {
	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	accessToken, err := generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Вы успешно вошли", "AccessToken для всех последующих операций": accessToken})

}

// @Summary Создание проекта
// @Description Создаёт новый проект, привязывая его к пользователю, авторизованному через JWT-токен
// @Tags Проекты
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Param input body Project true "Данные проекта"
// @Success 201 {object} map[string]interface{} "Проект успешно создан"
// @Failure 400 {object} map[string]string "Некорректный ввод данных"
// @Failure 401 {object} map[string]string "Необходим авторизационный токен или неверный формат токена"
// @Failure 500 {object} map[string]string "Ошибка базы данных"
// @Router /projects [post]
func createProject(c *gin.Context) {
	// Проверяем подключение к базе
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем токен из заголовка
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	// Убираем "Bearer " перед токеном
	const prefix = "Bearer "
	if len(tokenString) > len(prefix) && strings.HasPrefix(tokenString, prefix) {
		tokenString = tokenString[len(prefix):]
	}

	// Разбираем JWT-токен
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Получаем ID пользователя из токена
	userID := claims.UserID

	// Читаем JSON-запрос
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Валидация полей проекта
	if project.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project name is required"})
		return
	}

	if len(project.Name) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project name must be at least 2 characters long"})
		return
	}

	if project.Description != "" && len(project.Description) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project description cannot exceed 500 characters"})
		return
	}

	// Присваиваем assignee_id автоматически
	project.AssigneeID = userID

	// Устанавливаем контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	// Создаём проект в базе
	if err := db.WithContext(ctx).Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusCreated, gin.H{"message": "Проект успешно создан", "Projects:": project})
}

// @Summary Получение проекта
// @Description Возвращает данные проекта по его ID
// @Tags Проекты
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param Authorization header string true "Bearer токен"
// @Success 200 {object} map[string]interface{} "Проект успешно найден"
// @Failure 400 {object} map[string]string "Некорректный ID проекта"
// @Failure 404 {object} map[string]string "Проект не найден"
// @Failure 500 {object} map[string]string "Ошибка базы данных"
// @Router /projects/{id} [get]
func getProject(c *gin.Context) {
	// Проверяем подключение к базе
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем ID из параметра URL и конвертируем в число
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil || projectID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Устанавливаем контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	// Ищем проект в базе данных
	var project Project
	err = db.WithContext(ctx).First(&project, projectID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}

	// Возвращаем найденный проект
	c.JSON(http.StatusOK, gin.H{"message": "Проект успешно найден", "Projects:": project})
}

// @Summary Получение проектов пользователя
// @Description Возвращает список проектов, созданных текущим пользователем
// @Tags Проекты
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer токен"
// @Success 200 {object} map[string]interface{} "Проекты успешно найдены"
// @Failure 401 {object} map[string]string "Неавторизованный доступ или неверный токен"
// @Failure 404 {object} map[string]string "У пользователя нет проектов"
// @Failure 500 {object} map[string]string "Ошибка базы данных"
// @Router /user/projects [get]
func getUserProjects(c *gin.Context) {
	// Проверяем подключение к базе
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем токен из заголовка Authorization
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Убираем "Bearer " из строки токена
	const prefix = "Bearer "
	if len(tokenString) <= len(prefix) || !strings.HasPrefix(tokenString, prefix) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}
	tokenString = tokenString[len(prefix):]

	// Разбираем токен и извлекаем данные
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Получаем UserID из токена
	userID := claims.UserID
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing user ID"})
		return
	}

	// Устанавливаем контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Ищем проекты, созданные пользователем
	var projects []Project
	err = db.WithContext(ctx).Where("assignee_id = ?", userID).Find(&projects).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	// Если проектов нет
	if len(projects) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No projects found for this user"})
		return
	}

	// Возвращаем найденные проекты
	c.JSON(http.StatusOK, gin.H{"message": "Проекты успешно найдены", "Projects:": projects})
}

// @Summary Обновление проекта
// @Description Обновляет данные существующего проекта по его ID
// @Tags Проекты
// @Accept json
// @Produce json
// @Param id path int true "ID проекта"
// @Param Authorization header string true "Bearer токен"
// @Param project body Project true "Данные для обновления проекта"
// @Success 200 {object} map[string]interface{} "Проект успешно обновлён"
// @Failure 400 {object} map[string]string "Неверный запрос или некорректные данные"
// @Failure 401 {object} map[string]string "Неавторизованный доступ"
// @Failure 404 {object} map[string]string "Проект не найден"
// @Failure 500 {object} map[string]string "Ошибка базы данных"
// @Router /projects/{id} [put]
func updateProject(c *gin.Context) {
	// Проверяем подключение к базе
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем ID из параметра
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Ищем проект в базе
	var project Project
	if err := db.First(&project, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Парсим новый JSON, но не затираем ID
	var updatedData Project
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Валидация полей структуры Project
	if updatedData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project name is required"})
		return
	}

	if len(updatedData.Name) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project name must be at least 2 characters long"})
		return
	}

	if updatedData.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project description is required"})
		return
	}

	// Обновляем только переданные поля, ID оставляем старый
	db.Model(&project).Updates(updatedData)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Проверяем ошибки при обновлении
	if err := db.WithContext(ctx).Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project: " + err.Error()})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Проект успешно обновлён", "Projects:": project})
}

// @Summary Загрузка файла к проекту
// @Description Загружает файл для указанного проекта и сохраняет путь к файлу в базе данных
// @Tags Проекты
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID проекта"
// @Param Authorization header string true "Bearer токен"
// @Param file formData file true "Файл для загрузки (до 100MB)"
// @Success 200 {object} map[string]interface{} "Файл успешно загружен"
// @Failure 400 {object} map[string]string "Некорректный запрос или ошибка загрузки файла"
// @Failure 401 {object} map[string]string "Неавторизованный доступ"
// @Failure 404 {object} map[string]string "Проект не найден"
// @Failure 500 {object} map[string]string "Ошибка при сохранении файла или обновлении базы данных"
// @Router /projects/{id}/upload [post]
func uploadProjectFile(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("id")) //(c.Param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	// Валидация размера файла (например, ограничиваем до 10 MB)
	const maxFileSize = 100 * 1024 * 1024 // 10 MB
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the 100MB limit"})
		return
	}

	// Создаём папку uploads (если её нет)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Формируем путь к файлу
	filePath := filepath.Join(uploadDir, fmt.Sprintf("project_%d_%s", projectID, file.Filename))

	// Сохраняем файл
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Сохраняем путь в БД
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	if err := db.WithContext(ctx).Model(&Project{}).Where("id = ?", projectID).Update("file_path", filePath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database with file path"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно загружен", "File_path": filePath})
}

// @Summary Скачивание файла проекта
// @Description Скачивает файл, связанный с указанным проектом, если он существует на сервере
// @Tags Проекты
// @Produce octet-stream
// @Param id path int true "ID проекта"
// @Param Authorization header string true "Bearer токен"
// @Success 200 {file} string "Файл проекта"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 401 {object} map[string]string "Неавторизованный доступ"
// @Failure 404 {object} map[string]string "Файл не найден"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /projects/{id}/download [get]
func downloadProjectFile(c *gin.Context) {
	// Проверяем подключение к базе
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем ID проекта из параметра и проверяем, что это число
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Запрашиваем путь файла из БД
	var filePath string
	err = db.Raw("SELECT file_path FROM projects WHERE id = ?", projectID).Scan(&filePath).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Если файл не найден
	if filePath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Проверяем, существует ли файл по этому пути
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File does not exist on server"})
		return
	}

	// Отправляем файл клиенту
	c.File(filePath)
}

// @Summary Удаление проекта
// @Description Удаляет проект по указанному ID
// @Tags Проекты
// @Param id path int true "ID проекта"
// @Param Authorization header string true "Bearer токен"
// @Success 200 {object} map[string]string "Проект удален"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 401 {object} map[string]string "Неавторизованный доступ"
// @Failure 404 {object} map[string]string "Проект не найден"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /projects/{id} [delete]
func deleteProject(c *gin.Context) {
	// Проверяем подключение к базе
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем ID проекта и проверяем, что это число
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Ищем проект в базе
	var project Project
	if err := db.First(&project, projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Удаляем проект
	if err := db.WithContext(ctx).Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project: " + err.Error()})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Проект успешно удалён"})
}

// Управление задачами
// @Summary Создание задачи
// @Description Создает новую задачу для проекта с привязкой к пользователю
// @Tags Задачи
// @Param id path int true "ID проекта"
// @Param Authorization header string true "Bearer токен"
// @Param task body Task true "Данные задачи"
// @Success 201 {object} map[string]interface{} "Задача успешно создана"
// @Failure 400 {object} map[string]string "Ошибка валидации данных"
// @Failure 401 {object} map[string]string "Неавторизованный доступ"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /projects/{id}/tasks [post]
func createTask(c *gin.Context) {
	var task Task
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем ID проекта из URL
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Получаем токен из заголовка
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
		return
	}

	// Убираем "Bearer " перед токеном
	const prefix = "Bearer "
	if len(tokenString) > len(prefix) && strings.HasPrefix(tokenString, prefix) {
		tokenString = tokenString[len(prefix):]
	}

	// Разбираем JWT-токен
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Получаем ID пользователя из токена
	userID := claims.UserID

	// Привязываем JSON и проверяем обязательные поля
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Устанавливаем ProjectID из параметра URL
	task.ProjectID = uint(projectID)
	task.AssigneeID = userID

	// Валидация данных
	if err := validate.Struct(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Validation failed: %v", err.Error())})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	// Сохраняем в базе
	if err := db.WithContext(ctx).Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "Err": err.Error()})
		return
	}

	// Отправляем ответ с созданной задачей
	c.JSON(http.StatusCreated, gin.H{"message": "Задача успешно создана", "Task": task})
}

// @Summary Получение задач проекта
// @Description Получает список задач проекта с возможностью фильтрации по статусу, дедлайну и приоритету
// @Tags Задачи
// @Param id path int true "ID проекта"
// @Param Authorization header string true "Bearer токен"
// @Param status query string false "Статус задачи (In_Progress, Done, In_Line)"
// @Param deadline query string false "Дедлайн задачи (формат: YYYY-MM-DD)"
// @Param priority query string false "Приоритет задачи (High, Medium, Low)"
// @Success 200 {object} map[string]interface{} "Список задач"
// @Failure 400 {object} map[string]string "Ошибка валидации данных"
// @Failure 404 {object} map[string]string "Задачи не найдены"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /projects/{id}/tasks [get]
func getTasks(c *gin.Context) {
	// Проверяем подключение к базе
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем ID проекта и проверяем его корректность
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Получаем параметры фильтрации
	status := c.Query("status")
	deadline := c.Query("deadline")
	priority := c.Query("priority")

	// Валидация значений статуса и приоритета
	if status != "" && !isValidStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status, allowed values are: In_Progress, Done, In_Line"})
		return
	}

	if priority != "" && !isValidPriority(priority) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority, allowed values are: High, Medium, Low"})
		return
	}

	// Валидация формата даты для deadline
	if deadline != "" {
		if _, err := time.Parse("2006-01-02", deadline); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format, expected YYYY-MM-DD"})
			return
		}
	}

	// Инициализируем запрос
	var tasks []Task
	query := db.Where("project_id = ?", projectID)

	// Применяем фильтрацию по статусу, если указано
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Применяем фильтрацию по deadline, если указано
	if deadline != "" {
		query = query.Where("date_trunc('day', deadline) = ?", deadline)
	}

	// Применяем фильтрацию по приоритету, если указано
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Выполняем запрос
	if err := query.WithContext(ctx).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	// Проверяем, есть ли задачи
	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No tasks found for the given criteria"})
		return
	}

	// Возвращаем результат
	c.JSON(http.StatusOK, gin.H{"Задачи": tasks})
}

// @Summary Обновление задачи
// @Description Обновляет задачу в проекте по ID, с проверкой обязательных полей и значений
// @Tags Задачи
// @Param task_id path int true "ID задачи"
// @Param Authorization header string true "Bearer токен"
// @Param task body Task true "Обновленные данные задачи"
// @Success 200 {object} map[string]interface{} "Информация об обновленной задаче"
// @Failure 400 {object} map[string]string "Ошибка валидации данных"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /projects/:id/tasks/:task_id [put]
func updateTask(c *gin.Context) {
	// Проверяем подключение к базе
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем ID задачи из параметра "task_id"
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Ищем задачу в базе данных по taskID
	var task Task
	if err := db.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Привязываем данные из JSON
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация обязательных полей и значений
	if task.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	if !isValidStatus(task.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status, allowed values are: In_Progress, Done, In_Line"})
		return
	}

	if !isValidPriority(task.Priority) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority, allowed values are: High, Medium, Low"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Обновляем задачу в базе данных
	if err := db.WithContext(ctx).Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task: " + err.Error()})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Задача успешно обновлена", "Task": task})
}

// @Summary Удаление задачи
// @Description Удаляет задачу из базы данных по ID
// @Tags Задачи
// @Param task_id path int true "ID задачи"
// @Param Authorization header string true "Bearer токен"
// @Success 200 {object} map[string]string "Сообщение об успешном удалении задачи"
// @Failure 400 {object} map[string]string "Ошибка, если ID задачи некорректен"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /tasks/{task_id} [delete]
func deleteTask(c *gin.Context) {
	// Проверяем подключение к базе данных
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Получаем ID задачи и проверяем, является ли он числом
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID, must be a number"})
		return
	}

	// Проверяем, существует ли задача в базе данных
	var task Task
	if err := db.First(&task, taskID).Error; err != nil {
		// Если задача не найдена, возвращаем ошибку
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			// Если произошла другая ошибка при поиске задачи
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		}
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Удаляем задачу
	if err := db.WithContext(ctx).Delete(&task).Error; err != nil {
		// Если возникла ошибка при удалении
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task: " + err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Задача успешно удалена"})
}

// Проверка допустимых значений для status
func isValidStatus(status string) bool {
	allowedStatuses := []string{"In_Progress", "Done", "In_Line"}
	for _, s := range allowedStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// Проверка допустимых значений для priority
func isValidPriority(priority string) bool {
	allowedPriorities := []string{"High", "Medium", "Low"}
	for _, p := range allowedPriorities {
		if p == priority {
			return true
		}
	}
	return false
}
