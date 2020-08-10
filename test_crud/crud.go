package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type employee struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type responseStd struct {
	Message string `json:"message"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	if dbUser == "" || dbPass == "" || dbName == "" {
		panic("database configuration not set properly")
	}

	descriptor := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName)
	db, err := sql.Open(dbDriver, descriptor)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetList(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT * FROM employee ")
	if err != nil {
		panic(err.Error())
	}
	var (
		emp     employee
		empList []employee
	)

	for selDB.Next() {
		var id int
		var name, phone string
		err = selDB.Scan(&id, &name, &phone)
		if err != nil {
			panic(err.Error())
		}
		emp.ID = id
		emp.Name = name
		emp.Phone = phone
		empList = append(empList, emp)
	}

	defer db.Close()

	Response(w, r, empList)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM employee WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	emp := employee{}
	for selDB.Next() {
		var id int
		var name, phone string
		err = selDB.Scan(&id, &name, &phone)
		if err != nil {
			panic(err.Error())
		}
		emp.ID = id
		emp.Name = name
		emp.Phone = phone
	}

	defer db.Close()

	Response(w, r, emp)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var resp responseStd

	if r.Method == "POST" {
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		id := rand.Intn(9999)
		insForm, err := db.Prepare("INSERT INTO employee(id, name, phone) VALUES(?,?,?)")
		if err != nil {
			log.Println("Failed to insert ", err)
			resp.Message = "Failed to insert"
			Response(w, r, resp)
			return
		}
		_, err = insForm.Exec(id, name, phone)
		if err != nil {
			log.Println("Failed to insert", err)
			resp.Message = "Failed to insert"
			Response(w, r, resp)
			return
		}
		log.Println("INSERT: Name: " + name + " | Phone: " + phone)
		resp.Message = "Success to insert " + name
		Response(w, r, resp)
		return
	}
	defer db.Close()

}

func Update(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var resp responseStd
	if r.Method == "PUT" {
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		id := r.FormValue("id")
		insForm, err := db.Prepare("UPDATE employee SET name=?, phone=? WHERE id=?")
		defer db.Close()
		if err != nil {
			log.Println("Failed to update ", err)
			resp.Message = "Failed to update"
			Response(w, r, resp)
			return
		}
		_, err = insForm.Exec(name, phone, id)
		if err != nil {
			log.Println("Failed to update ", err)
			resp.Message = "Failed to update"
			Response(w, r, resp)
			return
		}
		log.Println("UPDATE: Name: " + name + " | Phone: " + phone)
		resp.Message = "Success to update " + name
		Response(w, r, resp)
		return
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {
		db := dbConn()
		var resp responseStd
		emp := r.URL.Query().Get("id")
		delForm, err := db.Prepare("DELETE FROM employee WHERE id=?")
		if err != nil {
			log.Println("Failed to update ", err)
			resp.Message = "Failed to update"
			Response(w, r, resp)
			return
		}
		delForm.Exec(emp)
		defer db.Close()

		resp.Message = "Success deleted data"
		Response(w, r, resp)
		return

	}
}

func main() {
	log.Println("Server started on: http://localhost:8080")

	http.HandleFunc("/list", GetList)
	http.HandleFunc("/detail", Detail)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.ListenAndServe(":8080", nil)
}

func Response(w http.ResponseWriter, r *http.Request, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
