package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
}
var client *mongo.Client
var taskCollection *mongo.Collection

func main() {
    // Initialize MongoDB 
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }
    taskCollection = client.Database("my_api").Collection("add_data")
    r := mux.NewRouter()

    r.HandleFunc("/tasks", createTask).Methods("POST")
    r.HandleFunc("/tasks", getAllTasks).Methods("GET")
    r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")


    port := ":8080"
    fmt.Printf("Server listening on port %s\n", port)
    err = http.ListenAndServe(port, r)
    if err != nil {
        log.Fatal(err)
    }
}
// Create a new task
func createTask(w http.ResponseWriter, r *http.Request) {
    var task Task
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Insert into MongoDB collection
    _, err := taskCollection.InsertOne(context.TODO(), task)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(task)
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
    cursor, err := taskCollection.Find(context.TODO(), bson.M{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    var tasks []Task
    for cursor.Next(context.TODO()) {
        var task Task
        if err := cursor.Decode(&task); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        tasks = append(tasks, task)
    }
    json.NewEncoder(w).Encode(tasks)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    taskID := params["id"]
    _, err := taskCollection.DeleteOne(context.TODO(), bson.M{"_id": taskID})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println( "Deleted")

    w.WriteHeader(http.StatusNoContent)
}