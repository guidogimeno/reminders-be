package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guidogimeno/reminders-be.git/api"
)

const port = ":3000"

func main() {
	fmt.Println("Starting server...")

	db, err := sql.Open("mysql", "root:codeblocks@tcp(localhost:3306)/remindersdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	mysql := api.NewMySQLService(db)
	server := api.NewApiServer(mysql)

	log.Fatal(server.Start(port))
}
