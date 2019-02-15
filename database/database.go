package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		fmt.Println("Old database connect ready")
		return db
	}
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("cannot connect to database server", err)
	}
	fmt.Println("new create database connect already")
	return db
}

func InsertCustomer(name, email, status string) *sql.Row {
	return Conn().QueryRow("INSERT INTO customers(name,email,status) VALUES ($1,$2,$3) RETURNING id", name, email, status)

}

func SelectCustomerAll() (*sql.Stmt, error) {
	return Conn().Prepare("SELECT * FROM customers")
}

func SelectCustomerbyID(id int) (*sql.Stmt, error) {
	return Conn().Prepare("SELECT * FROM customers WHERE id=$1")
}

func UpdateCustomer(id int, name, email, status string) (*sql.Stmt, error) {
	return Conn().Prepare("UPDATE customers SET name=$2,email=$3,status=$4 WHERE id=$1")
}

func DeleteCustomer(id int) (sql.Result, error) {
	return Conn().Exec("DELETE FROM customers WHERE id = $1", id)
}
