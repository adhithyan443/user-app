package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

//Connection to Postgres using database/sql

func ConnectDatabase(){

	psqlInfo:="host=localhost port=5432 user=postgres password=root dbname=user_app sslmode=disable"

	var err error

	DB,err = sql.Open("postgres", psqlInfo)

	if err != nil{
		log.Fatal("Failed to open database connection: ",err)
	}

	err = DB.Ping()

	if err != nil{
		log.Fatal("Failed to ping database: ",err)
	}

	fmt.Println("connected to Postgres")
}