package main

import (
	//"encoding/json"
	//"net/http"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")   //User
	cfg.Passwd = os.Getenv("DB_PASS") //Pass

	router := gin.Default()
	router.GET("ruta", method)
	router.Run("localhost:3912")
}
func method(c *gin.Context) {}
