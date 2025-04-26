package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sync"
)

type Activities struct {
	mu sync.Mutex
	db *sql.DB
}

const file string = "users.db"

func main() {

}

func NewUser(userName string, PW string) (string, error) {

	db, err := sql.Open("sqlite3", "/Users/aidancoco/Desktop/projects/fishTankWebGame/db/users.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	exists, err := checkIfUserExists(userName, db)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("user already exists")
	}

	hashPW, err := hashPassword(PW)
	if err != nil {
		return "", err
	}

	err = addUser(db, userName, hashPW)
	if err != nil {
		return "", err
	}
	return "user created successfully", nil
}

func addUser(db *sql.DB, userName string, hashedPW string) error {
	stmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(userName, hashedPW)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	fmt.Println("User created:", userName)
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckLoginUser(userName string, enteredPassword string) bool {
	db, err := sql.Open("sqlite3", "/Users/aidancoco/Desktop/projects/fishTankWebGame/db/users.db")
	if err != nil {
		log.Fatal(err)
	}

	exists, err := checkIfUserExists(userName, db)
	if !exists {
		println("no user found with name", userName)
		return false
	}

	hashPW, err := getPW(db, userName)
	if err != nil {
		log.Fatal("error retrieving pw:", err)
	}

	id, err := getUserID(db, userName)
	if err != nil {
		log.Fatal(err)
	}
	err = updateLastLogin(db, id)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	return checkPasswordHash(enteredPassword, hashPW)
}

func getPW(db *sql.DB, userName string) (string, error) {
	stmt, err := db.Prepare("SELECT PASSWORD FROM users WHERE USERNAME = ?")
	if err != nil {
		return "", fmt.Errorf("prepare insert: %w", err)
	}

	defer stmt.Close()

	pw := stmt.QueryRow(userName)

	var stringPW string

	err = pw.Scan(&stringPW)

	if err != nil {
		return "", fmt.Errorf("error scanning returned row")
	}

	return stringPW, nil
}

func checkIfUserExists(userName string, db *sql.DB) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ? LIMIT 1)`
	err := db.QueryRow(query, userName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil

}

func updateLastLogin(db *sql.DB, userID int) error {
	query := `UPDATE users SET last_login = current_timestamp WHERE id = ?`
	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("erorr updating logintime for userid %d, %w", userID, err)
	}
	return nil
}
