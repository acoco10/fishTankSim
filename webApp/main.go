package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"webApp/db"

	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
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
	mu        = sync.Mutex{}
	loginChan = make(chan string)
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

	println(">>> Method received:", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	println("trying to save game")
	var s GameState

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading save file"+string(err.Error()), http.StatusInternalServerError)
		return
	}
	b := strings.ReplaceAll(string(body), "\\\"", "\"")
	b = strings.ReplaceAll(string(b), "\"[", "[")
	b = strings.ReplaceAll(string(b), "]\"", "]")
	println("save JSON:\n", b)

	err = json.Unmarshal([]byte(b), &s)
	if err != nil {
		http.Error(w, "Error marshalling json for save"+string(err.Error()), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		http.Error(w, "Could not encode JSON"+string(err.Error()), http.StatusInternalServerError)
		return
	}

	err = db.Save(s.Username, string(data))
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write([]byte("Game saved"))
	if err != nil {
		log.Fatal(err)
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
	json.NewDecoder(r.Body).Decode(&u)

	dbSave, err := db.LoadSave(u.Username)

	if err != nil {
		http.Error(w, "Could not load save from db:"+string(err.Error()), http.StatusInternalServerError)
	}

	_, err = w.Write([]byte(dbSave))
	//w.Header().Set("Content-Type", "application/json")
	//err = json.NewEncoder(w).Encode(save)
	if err != nil {
		http.Error(w, "Could write save to json from db:"+string(err.Error()), http.StatusInternalServerError)
	}
}

func testWasmLocallyInitServer() {
	// You can now embed this URL in your HTML

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/load", loadHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/main.wasm", func(w http.ResponseWriter, r *http.Request) {
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

	if arg == "local" {
		testWasmLocallyInitServer()
	} else {

		initServer()

		log.Printf("Listening on :%s...", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}

	}
}
