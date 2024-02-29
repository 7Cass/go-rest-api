package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type User struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Age int `json:"age"`
}

var users []User

func main() {
  http.HandleFunc("/users", handleGetUsers)
  http.HandleFunc("/users/create", handleCreateUser)
  http.HandleFunc("/users/delete", handleDeleteUser)
  http.HandleFunc("/users/update", handleUpdateUser)

  fmt.Println("Server listening on port 8080...")
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.WriteHeader(http.StatusMethodNotAllowed)
    return
  }

  var user User

  err := json.NewDecoder(r.Body).Decode(&user)

  // Define `id` of new user
  if len(users) == 0 {
    user.Id = 1
  } else {
    lastUserId := users[len(users) - 1].Id
    user.Id = lastUserId + 1
  }

  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  
  users = append(users, user)
  
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(users)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    w.WriteHeader(http.StatusMethodNotAllowed)
    return
  }

  w.Header().Set("Content-Type", "application/json")

  if len(users) == 0 {
    w.Write([]byte("[]"))
    return
  }

  json.NewEncoder(w).Encode(users)
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodDelete {
    w.WriteHeader(http.StatusMethodNotAllowed)
    return
  }

  // Extract user `id` from query string
  id := r.URL.Query().Get("id")
  if id == "" {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // Convert `id` to integer
  userId, err := strconv.Atoi(id)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // Find the user in users slice
  index := -1
  for i, user := range users {
    if user.Id == userId {
      index = i
      break
    }
  }
  
  // Handle if user is not found
  if index == -1 {
    w.WriteHeader(http.StatusNotFound)
    return
  }

  // Remove user from users slice
  users = append(users[:index], users[index + 1:]...)

  w.WriteHeader(http.StatusOK)
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPatch {
    w.WriteHeader(http.StatusMethodNotAllowed)
    return
  }

  // Extract user `id` from query string 
  // and convert to integer
  id := r.URL.Query().Get("id")
  userId, err := strconv.Atoi(id)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // Decode request body in new user
  var updatedUser map[string]interface{}
  err = json.NewDecoder(r.Body).Decode(&updatedUser)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  // Found user by `id` in users slice
  found := false
  for i, user := range users {
    if user.Id == userId {
      // Update only the fields provided in body
      for key, value := range updatedUser {
        switch key {
          case "name":
            users[i].Name = value.(string)
          case "age":
            users[i].Age = int(value.(float64))
          default:
            w.WriteHeader(http.StatusBadRequest)
            return
        }
      }
      found = true
      // Return updated user
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(users[i])
      return

    }
  }

  if !found {
    w.WriteHeader(http.StatusNotFound)
    return
  }
}
