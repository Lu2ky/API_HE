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
type OfficialSchedule struct {
	N_iduser               int             `json:"N_iduser"`
	N_idcourse             int             `json:"N_idcourse"`
	Nrc                    string          `json:"Nrc"`
	Course                 string          `json:"Course"`
	Tag                    string          `json:"Tag"`
	Day                    int             `json:"Day"`
	StartHour              string          `json:"StartHour"`
	EndHour                string          `json:"EndHour"`
	Classroom              string          `json:"Classroom"`
	Credits                sql.NullFloat64 `json:"Credits"`
	Standardofcalification string          `json:"Standardofcalification"`
	Campus                 string          `json:"Campus"`
}
type PersonalSchedule struct {
	N_iduser    int    `json:"N_iduser"`
	N_idcourse  int    `json:"N_idcourse"`
	Activity    string `json:"Activity"`
	Tag         string `json:"Tag"`
	Description sql.NullString `json:"Description"`
	Dt_Start	sql.NullString `json:"Dt_Start"`
	Dt_End		sql.NullString `json:"Dt_End"`
	Day         int    `json:"Day"`
	StartHour   string `json:"StartHour"`
	EndHour     string `json:"EndHour"`
	IsDeleted   *sql.NullBool   `json:"IsDeleted"`
}
type PersonalScheduleNewValue struct {
	NewActivityValue   string `json:"NewActivityValue" binding:"required"`
	IdPersonalSchedule int    `json:"IdPersonalSchedule" binding:"required"`
}
type forDeleteOrRecoveryPersonalSchedule struct {
	IsDeleted          *bool `json:"IsDeleted" binding:"required"`
	IdPersonalSchedule int   `json:"IdPersonalSchedule" binding:"required"`
}
type NewPersonalActivity struct {
	Activity string `json:"Activity"`
	Description string `json:"Description"`
	IdTag       int `json:"IdTag"`
	Day         int    `json:"Day"`
	StartHour   string `json:"StartHour"`
	EndHour     string `json:"EndHour"`
	N_iduser    int    `json:"N_iduser"`
	Id_AcademicPeriod int `json:"Id_AcademicPeriod"`
}
func apiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		validAPIKey := os.Getenv("API_KEY")
		if validAPIKey == "" {
			log.Fatal("API_KEY no configurada en el archivo .env")
		}
		if apiKey == "" {
			c.JSON(401, gin.H{"error": "API Key requerida. Incluya el header X-API-Key"})
			c.Abort()
			return
		}
		if apiKey != validAPIKey {
			c.JSON(403, gin.H{"error": "API Key inválida"})
			c.Abort()
			return
		}
		
		c.Next()
	}
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
	router := gin.Default()                     //Create the default router for POST/GET methods
	//router.GET("/GetUserById/:id", getUserById) /* Use the / for subdirectorys in the localhost:3912 and references the method */
	router.GET("/GetOfficialScheduleByUserId/:id", getOfficialScheduleByUserId)
	router.GET("/GetPersonalScheduleByUserId/:id", getPersonalScheduleByUserId)
	router.POST("/updateNameOfPersonalScheduleByIdCourse", updateNameOfPersonalScheduleByIdCourse)
	router.POST("/updateDescriptionOfPersonalScheduleByIdCourse", updateDescriptionOfPersonalScheduleByIdCourse)
	router.POST("/updateStartHourOfPersonalScheduleByIdCourse", updateStartHourOfPersonalScheduleByIdCourse)
	router.POST("/updateEndHourOfPersonalScheduleByIdCourse", updateEndHourOfPersonalScheduleByIdCourse)
	router.POST("/deleteOrRecoveryPersonalScheduleByIdCourse", deleteOrRecoveryPersonalScheduleByIdCourse)
	router.POST("/addPersonalActivity", addPersonalActivity)
	router.Run("0.0.0.0:3913") // The port number for expone the API
}
func method(c *gin.Context) {}

// c *gin.Context essential for method in GET/POST actions

/* This function is a basic get for get the users from database */


func getOfficialScheduleByUserId(c *gin.Context) {
	id := c.Param("id")

	rows, err := db.Query(`
		SELECT ao.*
		FROM ActividadesOficiales ao
		JOIN Usuarios u ON ao.N_idUsuario = u.N_idUsuario
		WHERE u.T_codUsuario = ?
	`, id)

	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	var ofcschedules []OfficialSchedule

	for rows.Next() {
		var ofcschedule OfficialSchedule
		err := rows.Scan(
			&ofcschedule.N_iduser,
			&ofcschedule.N_idcourse,
			&ofcschedule.Nrc,
			&ofcschedule.Course,
			&ofcschedule.Tag,
			&ofcschedule.Day,
			&ofcschedule.StartHour,
			&ofcschedule.EndHour,
			&ofcschedule.Classroom,
			&ofcschedule.Credits,
			&ofcschedule.Standardofcalification,
			&ofcschedule.Campus,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			c.JSON(500, gin.H{"error": "Error en procesamiento de datos"})
			return
		}
		ofcschedules = append(ofcschedules, ofcschedule)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		c.JSON(500, gin.H{"error": "Error leyendo resultados"})
		return
	}

	c.JSON(200, ofcschedules)
}

func getPersonalScheduleByUserId(c *gin.Context) {
	id := c.Param("id")
	var rows *sql.Rows
		rows, err := db.Query(`
		SELECT ao.*
		FROM ActividadesPersonales ao
		JOIN Usuarios u ON ao.N_idUsuario = u.N_idUsuario
		WHERE u.T_codUsuario = ?
	`, id)

	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	defer rows.Close()

	var perschedules []PersonalSchedule
	for rows.Next() {
		var perschedule PersonalSchedule
		err := rows.Scan(&perschedule.N_iduser, 
			&perschedule.N_idcourse, 
			&perschedule.Activity, &perschedule.Tag, 
			&perschedule.Description, 
			&perschedule.Dt_Start,
			&perschedule.Dt_End,
			&perschedule.Day, 
			&perschedule.StartHour, 
			&perschedule.EndHour, 
			&perschedule.IsDeleted)
		if err != nil {
			log.Printf("Scan error: %v", err)
			c.JSON(500, gin.H{"Error": "Error en procesamiento de datos"})
			return
		}
		perschedules = append(perschedules, perschedule)

	}
	c.JSON(200, perschedules)
}
func updateNameOfPersonalScheduleByIdCourse(c *gin.Context) {
	var newValue PersonalScheduleNewValue
	err := c.BindJSON(&newValue)
	if err != nil {
		c.JSON(400, gin.H{"Palurdo": "formato invalido de json"})
		return
	}
	result, err := db.Exec("UPDATE ActividadesPersonales SET Actividad = ? WHERE N_idCurso= ? ", newValue.NewActivityValue, newValue.IdPersonalSchedule)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Personal schedule not found"})
		return
	}

	c.JSON(200, gin.H{
		"message":      "Personal schedule updated successfully",
		"rowsAffected": rowsAffected,
	})
}
func updateDescriptionOfPersonalScheduleByIdCourse(c *gin.Context) {
	var newValue PersonalScheduleNewValue
	err := c.BindJSON(&newValue)
	if err != nil {
		c.JSON(400, gin.H{"Palurdo": "formato invalido de json"})
		return
	}
	result, err := db.Exec("UPDATE ActividadesPersonales SET Descripcion = ? WHERE N_idCurso= ? ", newValue.NewActivityValue, newValue.IdPersonalSchedule)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Personal schedule not found"})
		return
	}

	c.JSON(200, gin.H{
		"message":      "Personal schedule updated successfully",
		"rowsAffected": rowsAffected,
	})
}
func updateStartHourOfPersonalScheduleByIdCourse(c *gin.Context) {
	var newValue PersonalScheduleNewValue
	err := c.BindJSON(&newValue)
	if err != nil {
		c.JSON(400, gin.H{"Palurdo": "formato invalido de json"})
		return
	}
	result, err := db.Exec("UPDATE ActividadesPersonales SET Hora_Inicio = ? WHERE N_idCurso= ? ", newValue.NewActivityValue, newValue.IdPersonalSchedule)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Personal schedule not found"})
		return
	}

	c.JSON(200, gin.H{
		"message":      "Personal schedule updated successfully",
		"rowsAffected": rowsAffected,
	})
}
func updateEndHourOfPersonalScheduleByIdCourse(c *gin.Context) {
	var newValue PersonalScheduleNewValue
	err := c.BindJSON(&newValue)
	if err != nil {
		c.JSON(400, gin.H{"Palurdo": "formato invalido de json"})
		return
	}
	result, err := db.Exec("UPDATE ActividadesPersonales SET Hora_Fin = ? WHERE N_idCurso= ? ", newValue.NewActivityValue, newValue.IdPersonalSchedule)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Personal schedule not found"})
		return
	}

	c.JSON(200, gin.H{
		"message":      "Personal schedule updated successfully",
		"rowsAffected": rowsAffected,
	})
}
func deleteOrRecoveryPersonalScheduleByIdCourse(c *gin.Context) {
	var deleteValue forDeleteOrRecoveryPersonalSchedule
	err := c.BindJSON(&deleteValue)
	if err != nil {
		c.JSON(400, gin.H{"Palurdo": "formato invalido de json"})
		return
	}
	result, err := db.Exec("UPDATE ActividadesPersonales SET B_isDeleted = ? WHERE N_idCurso=?", deleteValue.IsDeleted, deleteValue.IdPersonalSchedule)
	if err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Personal schedule not found"})
		return
	}

	c.JSON(200, gin.H{
		"message":      "Personal schedule updated successfully",
		"rowsAffected": rowsAffected,
	})
}
func addPersonalActivity(c *gin.Context) {
	var newPerActivity NewPersonalActivity
	err := c.BindJSON(&newPerActivity)
	if err != nil {
		c.JSON(400, gin.H{"error": "formato invalido de json"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Transaction error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	_, err0 := tx.Exec(
		"INSERT INTO Cursos (T_nombre, N_idEtiqueta, T_descripcion) VALUES (?, ?, ?);",
		newPerActivity.Activity,
		newPerActivity.IdTag,
		newPerActivity.Description,
	)
	if err0 != nil {
		tx.Rollback()
		log.Printf("Database error: %v", err0)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	_, err1 := tx.Exec(
		"INSERT INTO dias_clase(N_dia, TM_horaInicio, TM_horaFin) VALUES (?,?,?);",
		newPerActivity.Day,
		newPerActivity.StartHour,
		newPerActivity.EndHour,
	)
	if err1 != nil {
		tx.Rollback()
		log.Printf("Database error: %v", err1)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	_, err = tx.Exec(
		"INSERT INTO Materia_has_dias_clase (N_idCurso, N_idDiasClase) VALUES ((SELECT N_idCurso FROM Cursos WHERE T_nombre = ? AND T_descripcion = ?), (SELECT N_idDiasCase FROM dias_clase WHERE N_dia = ? AND TM_horaInicio = ? AND TM_horaFin = ?);",
		newPerActivity.Activity,
		newPerActivity.Description,
		newPerActivity.Day,
		newPerActivity.StartHour,
		newPerActivity.EndHour)
	if err != nil {
		tx.Rollback()
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	_, err = tx.Exec(
		"INSERT INTO horario (N_idUsuario, N_idCurso, N_idPeriodoAcademico) VALUES (?, (SELECT N_idCurso FROM Cursos WHERE T_nombre = ? AND T_descripcion = ?),?);",
		newPerActivity.N_iduser,
		newPerActivity.Activity,
		newPerActivity.Description,
		newPerActivity.Id_AcademicPeriod)
	if err != nil {
		tx.Rollback()
		log.Printf("Database error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Commit error: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{
		"message":    "Actividad creada correctamente",
	})
}

