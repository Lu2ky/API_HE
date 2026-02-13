package main

import (
	//"encoding/json"
	//"net/http"
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

/* Saving the session of MySQL, this is global for the access in all methods */

type User struct {
	Id       int    `json:"T_idUsuario"`
	Name     string `json:"T_nombre"`
	Programa string `json:"T_programa"`
}

func main() {
	err := godotenv.Load() // Load enviorement variables
	if err != nil {
		log.Fatal(".env file (error corrupted/not found)")
	}
	cfg := mysql.NewConfig()          //Create the cfg for MySQL
	cfg.User = os.Getenv("DB_USER")   //User
	cfg.Passwd = os.Getenv("DB_PASS") //Pass
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_ADDR") + ":" + os.Getenv("DB_ADDR_PORT")
	cfg.DBName = os.Getenv("DB_NAME")
	var err2 error
	db, err2 = sql.Open("mysql", cfg.FormatDSN())
	if err2 != nil {
		log.Fatal("Error connecting to database:", err2)
	}
	defer db.Close()
	defer db.Close()
	router := gin.Default()                     //Create the default router for POST/GET methods
	router.GET("/GetUserById/:id", getUserById) /* Use the / for subdirectorys in the localhost:3912
	and references the method 						*/
	router.Run("localhost:3912") // The port number for expone the API
}
func method(c *gin.Context) {}

// c *gin.Context essential for method in GET/POST actions

/* This function is a basic get for get the users from database */

func getUserById(c *gin.Context) {
	id := c.Param("id")
	/* Extract the id param for the query	*/
	var UserSaved User
	/* Put the param in the query remplace '?' for the id */
	err := db.QueryRow("SELECT T_codUsuario, T_nombre, T_programa FROM Usuarios WHERE T_codUsuario = ?", id).Scan(&UserSaved.Id, &UserSaved.Name, &UserSaved.Programa)
	/* test the insert the query in the database */
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"error": "User not found"}) //ERROR 404
			return
		}
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"}) //ERROR 500
		return
	}
	c.JSON(200, UserSaved) /* 200 = post the json in http */
}
