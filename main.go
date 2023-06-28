package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
  "os"
  "gorm.io/driver/mysql"
  "gorm.io/gorm"
	"log"
  "github.com/joho/godotenv"
	"strings"
	"time"
	base64 "encoding/base64"
)

var DB *gorm.DB

type Model struct {
  ID        uint           `gorm:"primaryKey"`
  CreatedAt time.Time
  UpdatedAt time.Time
  DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Book struct {
	Model
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func main() {
	LoadEnv()
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, DBNAME)
	fmt.Println(URL)
  db, err := gorm.Open(mysql.Open(URL))
	if err != nil {
			panic("Failed to connect to database!")

	}
	fmt.Println("Database connection established")
	err = db.AutoMigrate(&Book{})
	if err != nil {
		fmt.Println("Failed to Automigrate")
		return
	}
	DB = db

	router := gin.Default()
	router.POST("/books", basicAuth(), postBooks)
	router.Run("localhost:8080")  
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
			log.Fatalf("unable to load .env file")
	}
}

func basicAuth() gin.HandlerFunc {
	API_USER_NAME := os.Getenv("BASIC_AUTH_USER_NAME")
	API_PASSWORD := os.Getenv("BASIC_AUTH_PASSWORD")

  return func(c *gin.Context) {
    auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

    if len(auth) != 2 || auth[0] != "Basic" {
      respondWithError(401, "Unauthorized", c)
      return
    }
    payload, _ := base64.StdEncoding.DecodeString(auth[1])
    pair := strings.SplitN(string(payload), ":", 2)

    if (pair[0] != API_USER_NAME || pair[1] != API_PASSWORD)   {
      respondWithError(401, "Unauthorized", c)
      return 
    }

    c.Next()
  }
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}
	c.JSON(code, resp)
	c.Abort()
}

func postBooks(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	result := DB.Create(&newBook)
	c.IndentedJSON(http.StatusCreated, result)
}
