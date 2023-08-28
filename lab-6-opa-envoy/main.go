package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Employee struct {
	gorm.Model
	Name string
}

var db *gorm.DB
var err error

func main() {
	start := time.Now()
	db, err = gorm.Open("postgres", "sslmode=disable host=postgres port=5432 user=postgres dbname=postgres password=postgres")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Employee{})

	db.Delete(&Employee{})
	db.Create(&Employee{Name: "Alice"})
	db.Create(&Employee{Name: "Bob"})
	db.Create(&Employee{Name: "Charlie"})
	db.Create(&Employee{Name: "Dan"})
	db.Create(&Employee{Name: "Eve"})

	router := mux.NewRouter()

	router.Path("/").Handler(http.FileServer(http.Dir(".")))

	router.HandleFunc("/api/employees", ListEmployees).Methods("GET")
	router.HandleFunc("/api/employees/{id:[0-9]+}", GetEmployee).Methods("GET")
	router.HandleFunc("/api/employees", CreateEmployee).Methods("POST")
	router.HandleFunc("/api/employees/{id:[0-9]+}", EditEmployee).Methods("PUT")
	router.HandleFunc("/api/employees/{id:[0-9]+}", DeleteEmployee).Methods("DELETE")

	log.Println("Server starting. Start time:")
	log.Println(time.Since(start))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func ListEmployees(w http.ResponseWriter, r *http.Request) {
	var employee []Employee
	db.Find(&employee)
	if err := json.NewEncoder(w).Encode(employee); err != nil {
		logError(err)
	}
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeId := vars["id"]

	if !checkIfEmployeeExists(employeeId) {
		if err := json.NewEncoder(w).Encode("Employee Not Found!"); err != nil {
			logError(err)
		}
		return
	}

	var employee Employee
	db.Where("id = ?", employeeId).First(&employee)
	if err := json.NewEncoder(w).Encode(employee); err != nil {
		logError(err)
	}
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		logError(err)
	}

	db.Create(&employee)
	w.WriteHeader(204)
}

func EditEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeId := vars["id"]

	var employee Employee
	db.Where("id = ?", employeeId).First(&employee)
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		logError(err)
	}

	db.Save(&employee)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeId := vars["id"]
	db.Where("id = ?", employeeId).Delete(&Employee{})
}

func checkIfEmployeeExists(employeeId string) bool {
	var employee Employee
	db.First(&employee, employeeId)
	return employee.Name != ""
}

func logError(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
}
