package customer

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khanaphoz/finalexam/database"
)

func getCustomerHandler(c *gin.Context) {
	stmt, err := database.SelectCustomerAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	var items = []Customer{}
	for rows.Next() {
		cus := Customer{}
		err := rows.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		items = append(items, cus)
	}
	c.JSON(http.StatusOK, items)
}

func getCustomerbyIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	stmt, err := database.SelectCustomerbyID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	row := stmt.QueryRow(id)
	cus := Customer{}
	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, cus)
}

func createCustomerHandler(c *gin.Context) {
	var item Customer
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	row := database.InsertCustomer(item.Name, item.Email, item.Status)
	err := row.Scan(&item.ID)
	if err != nil {
		c.JSON(http.StatusAccepted, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusCreated, item)
}

func updateCustomerHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item := Customer{}
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	stmt, err := database.UpdateCustomer(id, item.Name, item.Email, item.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if _, err := stmt.Exec(id, item.Name, item.Email, item.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	item.ID = id
	c.JSON(http.StatusOK, item)
}

func deleteCustomerHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	//sqlStatement := `DELETE FROM customers WHERE id = $1`

	if _, err := database.DeleteCustomer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}
func loginMiddleware(c *gin.Context) {
	//log.Println("starting middleware")
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
	//log.Println("ending middleware")

}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(loginMiddleware)
	v1 := r.Group("")
	v1.GET("/customers", getCustomerHandler)
	v1.GET("/customers/:id", getCustomerbyIDHandler)
	v1.POST("/customers", createCustomerHandler)
	v1.PUT("/customers/:id", updateCustomerHandler)
	v1.DELETE("/customers/:id", deleteCustomerHandler)
	return r
}

func CreateTable() {
	createTb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err := database.Conn().Exec(createTb)
	if err != nil {
		log.Fatal("cannot create table", err)
	}
	fmt.Println("Create table success")
}
