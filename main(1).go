package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func setUpDB(*sql.DB) (*sql.DB, error) {

	const (
		host     = "localhost"
		port     = "5432"
		user     = "admin"
		password = "password"
		dbname   = "users"
	)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	//	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

type application struct {
	db *sql.DB
}

func main() {

	var db, err = setUpDB(nil)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &application{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/user", app.createUserForm)
	mux.HandleFunc("/user-add", app.createUser)
	mux.HandleFunc("/display", app.displayListings)
	log.Println("Starting server on port :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
