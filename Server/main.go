package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GameState struct {
	Username string        `json:"username"`
	State    []interface{} `json:"state"`
}

var (
	users       = make(map[string]string) // username -> password
	gameSaves   = make(map[string]string) // username -> saved state
	mu          = sync.Mutex{}
	currentUser = "" // protects maps
	loginChan   = make(chan string)
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Only POST allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	stored, ok := users[u.Username]
	if !ok || stored != u.Password {
		http.Error(w, `{"error":"Invalid login"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	if err != nil {
		return
	}
	select {
	case loginChan <- u.Username:
	default: // do nothing if already sent
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	json.NewDecoder(r.Body).Decode(&u)

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[u.Username]; exists {
		http.Error(w, "User exists", http.StatusBadRequest)
		return
	}

	users[u.Username] = u.Password
	err := json.NewEncoder(w).Encode(map[string]string{"message": "Registration Successful"})
	gameSaves[u.Username] = "0"
	if err != nil {
		return
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	println("trying to save game")
	var s GameState

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	b := strings.ReplaceAll(string(body), "\\\"", "\"")
	b = strings.ReplaceAll(string(b), "\"[", "[")
	b = strings.ReplaceAll(string(b), "]\"", "]")
	println(b)

	err = json.Unmarshal([]byte(b), &s)
	if err != nil {
		log.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		http.Error(w, "Could not encode JSON", http.StatusInternalServerError)
		return
	}

	fileName := fmt.Sprintf("save%s.json", s.Username)

	err = os.WriteFile("../assets/data/"+fileName, data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write([]byte("Game saved"))
	if err != nil {
		log.Fatal(err)
	}
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only get allowed", http.StatusMethodNotAllowed)
		return
	}
	println("checking for save")
	//username := r.URL.Query().Get("user")
	save, err := os.ReadFile("../assets/data/saveaidan.json")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(save)
	if err != nil {
		log.Fatal(err)
	}

}

func initServer() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/main.wasm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		http.ServeFile(w, r, "static/main.wasm")
	})
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	initServer()
}
