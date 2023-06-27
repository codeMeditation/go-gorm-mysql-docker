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
)

var DB *gorm.DB

type Book struct {
	ID     string  `json:"id"`
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
	router.POST("/books", postBooks)
	router.Run("localhost:8080")  
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
			log.Fatalf("unable to load .env file")
	}
}

func postBooks(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	result := DB.Create(&newBook)
	c.IndentedJSON(http.StatusCreated, result)
}
