package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=Ek16Scl03 dbname=postgres sslmode=disable") // подключение БД
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	InitDB()

	r := mux.NewRouter()

	r.HandleFunc("/product/", GetProducts).Methods("GET")
	r.HandleFunc("/product/", CreateProduct).Methods("POST")
	r.HandleFunc("/product/{id}", GetProductByID).Methods("GET")
	r.HandleFunc("/product/{id}", DeleteProduct).Methods("DELETE")
	r.HandleFunc("/product/{id}", EditProduct).Methods("PUT")
	r.HandleFunc("/measure/", GetMeasure).Methods("GET")
	r.HandleFunc("/measure/{id}", GetMeasureByID).Methods("GET")
	r.HandleFunc("/measure/", CreateMeasure).Methods("POST")
	r.HandleFunc("/measure/{id}", EditMeasure).Methods("PUT")
	r.HandleFunc("/measure/{id}", DeleteMeasure).Methods("DELETE")

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
