package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetMeasure(w http.ResponseWriter, r *http.Request) {
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

func GetMeasureByID(w http.ResponseWriter, r *http.Request) {
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

func CreateMeasure(w http.ResponseWriter, r *http.Request) {
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

func EditMeasure(w http.ResponseWriter, r *http.Request) {
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

func DeleteMeasure(w http.ResponseWriter, r *http.Request) {
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
