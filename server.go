package main

import (
	"github.com/khanaphoz/finalexam/customer"
	"github.com/khanaphoz/finalexam/database"
)

func main() {
	database.Conn()
	customer.CreateTable()

	r := customer.NewRouter()

	r.Run(":2019")

}
