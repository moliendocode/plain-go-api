package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Name string
	email string
	ID string
	password string
}

type usersHandler struct {
	sync.Mutex
	Users	map[string]User
}

func (h *usersHandler) method(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			h.get(w, r)
			return
		case "POST":
			h.post(w, r)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
			return
		

	}
}

func (h *usersHandler) get(w http.ResponseWriter, r *http.Request) {
	users := make([]User, len(h.Users))

	h.Lock()
	i := 0
	for _, user := range h.Users {
		users[i] = user
		i++
	}
	h.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
	
}

func (h *usersHandler) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Content-Type must be application/json, but got '%s'", ct)))
		return
	}

	var user User
	json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.Users[user.ID] = user
	defer h.Unlock()


}

func newUsersHandler() *usersHandler {
	return &usersHandler{
		Users: map[string]User{
			"1": {
				Name: "John",
				email: "john@asd.com",
				ID: "1",
				password: "123",
		},
		"2": {
			Name: "Brandon",
			email: "brandon@asd.com",
			ID: "2",
			password: "321",
	},
},
}}

func main() {
	usersHandler := newUsersHandler()
	http.HandleFunc("/users", usersHandler.method)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}