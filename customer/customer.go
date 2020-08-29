package customer

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

type Response struct {
	Message string `json:"message"`
}

var db *sql.DB

func init() {

	var err error

	// pubConn to elephantDB
	pubConn := "postgres://qrulmoqq:3Zu_Xhw121TaveedaBFPEQ_5Z_MGVel6@john.db.elephantsql.com:5432/qrulmoqq"

	db, err = sql.Open("postgres", getEnv("DATABASE_URL", pubConn))

	if err != nil {
		log.Fatal(err)
	}

	_, err = createCustomersTbl()

	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func createCustomersTbl() (string, error) {

	customerTbl := `
		CREATE TABLE IF NOT EXISTS Customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
		);
	`
	_, err := db.Exec(customerTbl)

	if err != nil {
		return "can't create table", err
	}

	return "create table success", nil

}

func CreateCustomerHandler(c *gin.Context) {

	t := Customer{}

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	row := db.QueryRow("INSERT INTO Customers (name, email, status) values ($1, $2, $3) RETURNING id", t.Name, t.Email, t.Status)
	err := row.Scan(&t.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, t)

}

func GetCustomersHandler(c *gin.Context) {

	status := c.Query("status")

	stmt, err := db.Prepare("SELECT id, name, email, status FROM Customers")

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	Customers := []Customer{}
	for rows.Next() {
		t := Customer{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		Customers = append(Customers, t)
	}

	tt := []Customer{}

	for _, item := range Customers {
		if status != "" {
			if item.Status == status {
				tt = append(tt, item)
			}
		} else {
			tt = append(tt, item)
		}
	}

	c.JSON(http.StatusOK, tt)

}

func GetCustomerByIdHandler(c *gin.Context) {

	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, name, email, status FROM Customers where id=$1")

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)
	t := &Customer{}

	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, t)
}

func UpdateCustomersHandler(c *gin.Context) {

	id := c.Param("id")

	stmt, err := db.Prepare("SELECT id, name, email, status FROM Customers where id=$1")

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	row := stmt.QueryRow(id)
	t := &Customer{}

	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := c.ShouldBindJSON(t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err = db.Prepare("UPDATE Customers SET status=$2, name=$3, email=$4 WHERE id=$1;")

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if _, err := stmt.Exec(id, t.Status, t.Name, t.Email); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, t)
}

func DeleteCustomersHandler(c *gin.Context) {

	r := Response{}
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM Customers WHERE id = $1")

	if err != nil {
		log.Fatal("can't delete statement", err)
		r.Message = "can't delete customer"
	}

	if _, err := stmt.Exec(id); err != nil {
		log.Fatal("can't execute delete statment", err)
		r.Message = "can't delete customer"
	}

	r.Message = "customer deleted"

	c.JSON(http.StatusOK, r)

}
