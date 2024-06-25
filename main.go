package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	UnitCost float64 `json:"unit_cost"`
	Measure  int     `json:"measure"`
}

type Measure struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=Ek16Scl03 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT id, name, quantity, unit_cost, measure FROM products")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.UnitCost, &product.Measure)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var product Product
	err := db.QueryRow("SELECT id, name, quantity, unit_cost, measure FROM products WHERE id = $1", id).Scan(&product.ID, &product.Name, &product.Quantity, &product.UnitCost, &product.Measure)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params Product
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO products (name, quantity, unit_cost, measure) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(params.Name, params.Quantity, params.UnitCost, params.Measure).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Product created successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
}

func editProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params Product
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE products SET name = $1, quantity = $2, unit_cost = $3, measure = $4 WHERE id = $5")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	_, err = stmt.Exec(params.Name, params.Quantity, params.UnitCost, params.Measure, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Product updated successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Product deleted successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
}

func getMeasure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT ID, name FROM measures")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var measures []Measure
	for rows.Next() {
		var measure Measure
		err := rows.Scan(&measure.ID, &measure.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		measures = append(measures, measure)
	}

	json.NewEncoder(w).Encode(measures)
}

func getMeasureByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	var measure Measure
	err := db.QueryRow("SELECT id, name FROM measures WHERE id = $1", id).Scan(&measure.ID, &measure.Name)
	if err != nil {
		http.Error(w, "Measure not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(measure)
}

func createMeasure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params Measure
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO measures (name) VALUES ($1) RETURNING id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(params.Name).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Measure created successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
}

func editMeasure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params Measure
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE measures SET name = $1  WHERE id = $2")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	vars := mux.Vars(r)
	id := vars["id"]
	_, err = stmt.Exec(params.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Measure updated successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
}

func deleteMeasure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.Exec("DELETE FROM measures WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Measure deleted successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	initDB()

	r := mux.NewRouter()

	r.HandleFunc("/product/", getProducts).Methods("GET")
	r.HandleFunc("/product/", createProduct).Methods("POST")
	r.HandleFunc("/product/{id}", getProductByID).Methods("GET")
	r.HandleFunc("/product/{id}", deleteProduct).Methods("DELETE")
	r.HandleFunc("/product/{id}", editProduct).Methods("PUT")
	r.HandleFunc("/measure/", getMeasure).Methods("GET")
	r.HandleFunc("/measure/{id}", getMeasureByID).Methods("GET")
	r.HandleFunc("/measure/", createMeasure).Methods("POST")
	r.HandleFunc("/measure/{id}", editMeasure).Methods("PUT")
	r.HandleFunc("/measure/{id}", deleteMeasure).Methods("DELETE")

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
