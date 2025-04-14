package main

import (
	"encoding/json"
	"fishTankWebGame/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net/http"
	"sync"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GameState struct {
	Username string `json:"username"`
	State    string `json:"state"`
}

var (
	users     = make(map[string]string) // username -> password
	gameSaves = make(map[string]string) // username -> saved state
	mu        = sync.Mutex{}            // protects maps
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
	if err != nil {
		return
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var s GameState
	json.NewDecoder(r.Body).Decode(&s)

	mu.Lock()
	defer mu.Unlock()

	gameSaves[s.Username] = s.State
	w.Write([]byte("Game saved"))
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	mu.Lock()
	state, ok := gameSaves[username]
	mu.Unlock()

	if !ok {
		http.Error(w, "No save found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(GameState{Username: username, State: state})
}

func initServer() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func runGame() {
	g := game.NewGame()
	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//initServer()
	runGame()
}
