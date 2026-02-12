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
	err := godotenv.Load() // Load enviorement variables
	if err != nil {
		log.Fatal(".env file (error corrupted/not found)")
	}
	cfg := mysql.NewConfig()          //Create the cfg for MySQL
	cfg.User = os.Getenv("DB_USER")   //User
	cfg.Passwd = os.Getenv("DB_PASS") //Pass

	router := gin.Default()    //Create the default router for POST/GET methods
	router.GET("ruta", method) /* Use the / for subdirectorys in the localhost:3912
	and references the method 						*/
	router.Run("localhost:3912") // The port number for expone the API
}
func method(c *gin.Context) {} // c *gin.Context essential for method in GET/POST actions
