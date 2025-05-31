package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/cors"
	"webApp/db"

	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SaveRequest struct {
	Username string    `json:"username"`
	Game     GameState `json:"game"`
}

type GameState struct {
	State []FishState `json:"state"`
	Tasks []Task      `json:"tasks"`
}

type FishState struct {
	Name      string `json:"Name"`
	Size      int    `json:"Size"`
	Progress  int    `json:"Progress"`
	NextLevel int    `json:"NextLevel"`
	FishType  string `json:"FishType"`
	MaxSpeed  int    `json:"MaxSpeed"`
}

type Task struct {
	Text      string `json:"Text"`
	Name      string `json:"Name"`
	Completed bool   `json:"Completed"`
}

var (
	mu        = sync.Mutex{}
	loginChan = make(chan string)
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Only POST allowed"}`+r.Method, http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	dbCheck, err := db.CheckLoginUser(u.Username, u.Password)
	if !dbCheck {
		http.Error(w, `{"error":"Invalid login"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
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
		log.Printf("Invalid method: %s at /register", r.Method)
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	log.Printf("Register attempt: username=%s", u.Username)

	mu.Lock()
	defer mu.Unlock()

	msg, err := db.NewUser(u.Username, u.Password)
	if err != nil {
		log.Printf("Failed to register user '%s': %v", u.Username, err)
		http.Error(w, `{"error":"Registration failed"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("User '%s' registered successfully", u.Username)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": msg}); err != nil {
		log.Printf("Failed to encode response for user '%s': %v", u.Username, err)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(">>> Method received:", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("trying to save game")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, "error reading save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Raw body:", string(body))

	var req SaveRequest

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("Unmarshal error:", err)
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	log.Println("Parsed State:", req.Game.State)
	log.Println("Parsed Tasks:", req.Game.Tasks)

	log.Println("Parsed GameState:", req.Game.State)

	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(req.Game, "", "  ")
	if err != nil {
		log.Println("Error re-encoding JSON:", err)
		http.Error(w, "could not encode JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Saving to DB with username:", req.Username)
	err = db.Save(req.Username, string(data))
	if err != nil {
		log.Println("Error saving to DB:", err)
		http.Error(w, "could not save to DB: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("Game saved"))
	if err != nil {
		log.Println("Could not write response:", err)
	}
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post allowed", http.StatusMethodNotAllowed)
		return
	}

	println("checking for save")
	//username := r.URL.Query().Get("user")

	/*	save, err := os.ReadFile("../assets/data/saveaidan.json")
		if err != nil {
			log.Fatal(err)
		}*/

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Could read user name from body:"+string(err.Error()), http.StatusInternalServerError)
	}

	dbSave, err := db.LoadSave(u.Username)

	if err != nil {
		http.Error(w, "Could not load save from db:"+string(err.Error()), http.StatusInternalServerError)
	}

	println(dbSave)

	_, err = w.Write([]byte(dbSave))
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, "Couldn't write save to json from db:"+string(err.Error()), http.StatusInternalServerError)
	}
}

func testWasmLocallyInitServer() {
	// You can now embed this URL in your HTML

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("static/fishFishFish/main.wasm", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received for /main.wasm")
		w.Header().Set("Content-Type", "application/wasm")
		http.ServeFile(w, r, "static/main.wasm")
	})
}

func generatePresignedURL() string {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	presigner := s3.NewPresignClient(client)

	req, err := presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("fish-fish-fish-assets"),
		Key:    aws.String("wasm/main.wasm"),
	}, s3.WithPresignExpires(15*time.Minute)) // Expires in 15 minutes

	if err != nil {
		panic("unable to presign request, " + err.Error())
	}

	return req.URL
}

func handleGetWasmURL(w http.ResponseWriter, r *http.Request) {
	// Configuration

	// Generate the pre signed URL using utility function
	wasmURL := generatePresignedURL()

	fmt.Println("Pre-signed URL:", wasmURL)
	// Prepare JSON response
	response := map[string]string{"url": wasmURL}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

func initServer() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)

	http.HandleFunc("/get-wasm-url", handleGetWasmURL)

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	println("running with os argument:", arg)

	// Use ServeMux for clean separation of routes
	mux := http.NewServeMux()

	// Register all your handlers on the mux
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/save", saveHandler)
	mux.HandleFunc("/load", loadHandler)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	// Add CORS middleware
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://collisionposition.netlify.app"}, // no trailing slash!
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(mux)

	if arg == "local" {
		log.Println("Running locally on http://localhost:5000")
		log.Fatal(http.ListenAndServe(":5000", handler))
	} else {
		log.Printf("Listening on :%s...", port)
		log.Fatal(http.ListenAndServe(":"+port, handler))
	}
}
