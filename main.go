package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"database/sql"
	"strconv"
    "github.com/rs/cors"
    "time"
    //"os"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)
type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Date      string    `json:"date"`
	Done      bool      `json:"done"`
	CreatedAt string `json:"created_at"`
}



func main() {
    //var DATABASEUSER = os.Getenv("DATABASEUSER")
    //var DATABASEPASSWORD = os.Getenv("DATABASEPASSWORD")
    // var DATABASEUSER = "app_db"
    // var DATABASEPASSWORD = "app_pass+123"
    // var DATABASENAME = "app_db"
    // var DATABASEHOST = "db"
    // var DATABASEPORT = "3306"
    //var EXPOSEDPORT = "8001"

    // fmt.Println("DATABASEUSER: ", DATABASEUSER)
    // fmt.Println("DATABASEPASSWORD: ", DATABASEPASSWORD)
    // fmt.Println("DATABASENAME: ", DATABASENAME)
    // fmt.Println("DATABASEHOST: ", DATABASEHOST)
    // fmt.Println("DATABASEPORT: ", DATABASEPORT)
    // fmt.Println("EXPOSEDPORT: ", EXPOSEDPORT)



    db, err := sql.Open("mysql", "app_user:app_pass+123@tcp(db:3306)/app_db?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	maxRetries := 10
    retryInterval := 5 * time.Second

    for i := 0; i < maxRetries; i++ {
        _, err = db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), date VARCHAR(255), done BOOLEAN, created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
        if err == nil {
            break
        }

        fmt.Printf("Failed to connect to MySQL: %v. Retrying in %v...\n", err, retryInterval)
        time.Sleep(retryInterval)
    }

    if err != nil {
        panic(err.Error())
    }
    
    
	r := mux.NewRouter()

    // API routes
    r.HandleFunc("/tasks", getTasks).Methods(http.MethodGet, http.MethodOptions)
    r.HandleFunc("/tasks", writeTasks).Methods(http.MethodPost, http.MethodOptions)
    r.HandleFunc("/tasks/{taskId}", deleteTasks).Methods(http.MethodDelete, http.MethodOptions)
    r.HandleFunc("/tasks/{taskId}", completeTasks).Methods(http.MethodPut, http.MethodOptions)

    // Serve frontend files
    r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./frontend"))))

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
    })

    handler := c.Handler(r)

    port := 8001
    fmt.Printf("Server is running on http://localhost:%d\n", port)
    http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}


func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var tasks []Task
    db, err := sql.Open("mysql", "app_user:app_pass+123@tcp(db:3306)/app_db?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	results, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var task Task
		err = results.Scan(&task.ID, &task.Name, &task.Date, &task.Done, &task.CreatedAt)
		if err != nil {
			panic(err.Error())
		}
		tasks = append(tasks, task)
	}
	json.NewEncoder(w).Encode(tasks)
}

func writeTasks(w http.ResponseWriter, r *http.Request) {
	var newTask Task
    err := json.NewDecoder(r.Body).Decode(&newTask)
    if err != nil {
        fmt.Println("Error decoding JSON:", err)
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    fmt.Println("Received task:", newTask)
    db, err := sql.Open("mysql", "app_user:app_pass+123@tcp(db:3306)/app_db?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("INSERT INTO tasks (name, date, done) VALUES (?, ?, ?)", newTask.Name, newTask.Date, newTask.Done)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	json.NewEncoder(w).Encode(newTask)

}

func deleteTasks(w http.ResponseWriter, r *http.Request) {
    // Parse task ID from the request parameters
    vars := mux.Vars(r)
    taskIdStr, ok := vars["taskId"]
    if !ok {
        http.Error(w, "Task ID not provided", http.StatusBadRequest)
        return
    }

    // Convert task ID to integer
    taskId, err := strconv.Atoi(taskIdStr)
    if err != nil {
        http.Error(w, "Invalid Task ID", http.StatusBadRequest)
        return
    }

    // Perform deletion in the database (assuming taskId is the primary key)
    db, err := sql.Open("mysql", "app_user:app_pass+123@tcp(db:3306)/app_db?parseTime=true")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("DELETE FROM tasks WHERE id = ?", taskId)
    if err != nil {
        fmt.Println("Error deleting task:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Respond with success
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.WriteHeader(http.StatusOK)
}

func completeTasks(w http.ResponseWriter, r *http.Request) {
    // Parse task ID from the request parameters
    vars := mux.Vars(r)
    taskIdStr, ok := vars["taskId"]
    if !ok {
        http.Error(w, "Task ID not provided", http.StatusBadRequest)
        return
    }

    // Convert task ID to integer
    taskId, err := strconv.Atoi(taskIdStr)
    if err != nil {
        http.Error(w, "Invalid Task ID", http.StatusBadRequest)
        return
    }

    // Parse task data from the request body
    var updatedTask Task
    err = json.NewDecoder(r.Body).Decode(&updatedTask)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Perform the update in the database (assuming taskId is the primary key)
    db, err := sql.Open("mysql", "app_user:app_pass+123@tcp(db:3306)/app_db?parseTime=true")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    _, err = db.Exec("UPDATE tasks SET done = ? WHERE id = ?", updatedTask.Done, taskId)
    if err != nil {
        fmt.Println("Error updating task:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Respond with the updated task data
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(updatedTask)
}
