package main

import (
	"context"
	"encoding/json"
	"fishTankWebGame/webApp/db"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

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

	if !db.CheckLoginUser(u.Username, u.Password) {
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

	msg, err := db.NewUser(u.Username, u.Password)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"message": msg})
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

	err = db.Save(s.Username, string(data))
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	_, err = w.Write([]byte(dbSave))
	//w.Header().Set("Content-Type", "application/json")
	//err = json.NewEncoder(w).Encode(save)
	if err != nil {
		log.Fatal(err)
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
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
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

	// Generate the presigned URL using utility function
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
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {

	var arg string

	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	println("running with os argument:", arg)

	if arg == "local" {
		testWasmLocallyInitServer()
	} else {
		initServer()
	}
}
