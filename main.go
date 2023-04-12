package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guidogimeno/reminders-be.git/api"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting server...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, dbPort, database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	mysql := api.NewMySQLService(db)
	server := api.NewApiServer(mysql)

	log.Fatal(server.Start(port))
}
