package main

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Task struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Duedate string `json:"duedate"`
}

func connect() (*sqlx.DB, error) {
	bin, err := ioutil.ReadFile("/run/secrets/db-password")
	if err != nil {
		return nil, err
	}
	return sqlx.Open("postgres", fmt.Sprintf("postgres://postgres:%s@db:5432/example?sslmode=disable", string(bin)))
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer db.Close()

	tasks := []Task{}
	err = db.Select(&tasks, "SELECT id, title, COALESCE(description, '') description, COALESCE(duedate, '') duedate FROM tasks ORDER BY id ASC")
	if err != nil {
		fmt.Println(err)
		 return
	}

	json.NewEncoder(w).Encode(tasks)
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var t Task
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// if id is specified, refuse to create a resource, because this is POST
	// TODO: query db whether the resource exists
	if (t.Id != "") {
		w.WriteHeader(204)
		return
	}

	query := `INSERT INTO tasks(title, description, duedate) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, t.Title, t.Description, t.Duedate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func putTasks(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var t Task
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
		UPDATE tasks SET title=$1, description=$2, duedate=$3 WHERE id=$4;
	`
	_, err = db.Exec(query, t.Title, t.Description, t.Duedate, t.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	query = `
		INSERT INTO tasks (title, description, duedate)
		SELECT $1, $2, $3
		WHERE NOT EXISTS (SELECT 1 FROM tasks WHERE id=$4);
	`
	_, err = db.Exec(query, t.Title, t.Description, t.Duedate, t.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func patchTasks(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var t Task
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	query := `
		UPDATE tasks SET title=$1, description=$2, duedate=$3 WHERE id=$4;
	`
	_, err = db.Exec(query, t.Title, t.Description, t.Duedate, t.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func deleteTasks(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		w.WriteHeader(500)
		return
	}

	var t Task
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `DELETE FROM tasks WHERE id = $1`
	_, err = db.Exec(query, t.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func main() {
	log.Print("Prepare db...")
	if err := prepare(); err != nil {
		log.Fatal(err)
	}

	log.Print("Listening 8000")
	r := mux.NewRouter()
	r.HandleFunc("/", getTasks).Methods(http.MethodGet)
	r.HandleFunc("/", postTasks).Methods(http.MethodPost)
	r.HandleFunc("/", putTasks).Methods(http.MethodPut)
	r.HandleFunc("/", patchTasks).Methods(http.MethodPatch)
	r.HandleFunc("/", deleteTasks).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r)))
}

func prepare() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	for i := 0; i < 60; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if _, err := db.Exec("DROP TABLE IF EXISTS tasks"); err != nil {
		return err
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY, 
			title VARCHAR NOT NULL, CHECK (title <> ''), 
			description VARCHAR, 
			duedate VARCHAR
		)
		`); err != nil {
		return err
	}

	// testing: put data into db
	for i := 0; i < 5; i++ {
		if _, err := db.Exec("INSERT INTO tasks (title) VALUES ($1);", fmt.Sprintf("task #%d", i)); err != nil {
			return err
		}
	}
	return nil
}
