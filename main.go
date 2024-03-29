// main.go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/WeebNinja/go-gorm-db/db"
	"github.com/WeebNinja/go-gorm-db/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read database configuration from environment variables
	dbType := os.Getenv("DB_TYPE")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Connect to the database
	database, err := db.ConnectDatabase(dbType, dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// AutoMigrate the database
	err = database.AutoMigrate(&models.Item{}, &models.Student{}, &models.Subject{}, &models.User{}, &models.Teacher{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// Create repositories for each model
	itemRepo := models.NewItemRepository(database)
	studentRepo := models.NewStudentRepository(database)
	subjectRepo := models.NewSubjectRepository(database)

	// Initialize Gin router
	r := gin.Default()

	// กำหนด cors (Cross-Origin Resource Sharing)
	r.Use(cors.New(cors.Config{
		// 3000 คือ port ที่ใช้งานใน frontend react
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Item routes
	r.GET("/items", itemRepo.GetItems)
	r.POST("/items", itemRepo.PostItem)
	r.GET("/items/:id", itemRepo.GetItem)
	r.PUT("/items/:id", itemRepo.UpdateItem)
	r.DELETE("/items/:id", itemRepo.DeleteItem)

	// Student routes
	r.GET("/students", studentRepo.GetStudents)
	r.POST("/students", studentRepo.CreateStudent)
	r.GET("/students/:id", studentRepo.GetStudent)
	r.PUT("/students/:id", studentRepo.UpdateStudent)
	r.DELETE("/students/:id", studentRepo.DeleteStudent)

	// Subject routes
	r.GET("/subjects", subjectRepo.GetSubjects)
	r.POST("/subjects", subjectRepo.CreateSubject)
	r.GET("/subjects/:id", subjectRepo.GetSubject)
	r.PUT("/subjects/:id", subjectRepo.UpdateSubject)
	r.DELETE("/subjects/:id", subjectRepo.DeleteSubject)

	// สร้างตัวแปร userRepo เพื่อเรียกใช้งาน UserRepository
	userRepo := models.NewUserRepository(database)

	// api /users จะเป็นการเรียกใช้งานฟังก์ชัน GetUsers ใน UserRepository
	r.GET("/users", userRepo.GetUsers)

	// api /users จะเป็นการเรียกใช้งานฟังก์ชัน PostUser ใน UserRepository
	r.POST("/users", userRepo.PostUser)

	// api /users/:email จะเป็นการเรียกใช้งานฟังก์ชัน GetUser ใน UserRepository
	// /users/abc@example จะเป็นการส่งค่า email ที่เป็นตัวอักษร abc@example ไปยังฟังก์ชัน GetUser ใน UserRepository
	r.GET("/users/:email", userRepo.GetUser)

	// api /users/:email จะเป็นการเรียกใช้งานฟังก์ชัน UpdateUser ใน UserRepository
	r.PUT("/users/:email", userRepo.UpdateUser)

	// api /users/:email จะเป็นการเรียกใช้งานฟังก์ชัน DeleteUser ใน UserRepository
	r.DELETE("/users/:email", userRepo.DeleteUser)

	// api /users/login จะเป็นการเรียกใช้งานฟังก์ชัน Login ใน UserRepository
	r.POST("/users/login", userRepo.Login)

	// Teacher routes
	teacherRepo := models.NewTeacherRepository(database)

	// Get all teachers
	r.GET("/teachers", teacherRepo.GetTeachers)

	// Create a new teacher
	r.POST("/teachers", teacherRepo.CreateTeacher)

	// Get a specific teacher by ID
	r.GET("/teachers/:id", teacherRepo.GetTeacher)

	// Update an existing teacher
	r.PUT("/teachers/:id", teacherRepo.UpdateTeacher)

	// Delete a teacher
	r.DELETE("/teachers/:id", teacherRepo.DeleteTeacher)

	// 404 route
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	})

	// Run the server
	if err := r.Run(":5000"); err != nil {
		log.Fatalf("Server is not running: %v", err)
	}
}
