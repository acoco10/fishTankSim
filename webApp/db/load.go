package db

import (
	"database/sql"
	"fmt"
	"os"
)

func LoadSave(userName string) (string, error) {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return " ", err
	}

	defer db.Close()

	id, err := getUserID(db, userName)
	if err != nil {
		return "", fmt.Errorf("couldnt load user id for loading save: %w", err)
	}

	save, err := load(db, id)

	if err != nil {
		return "", fmt.Errorf("couldnt load user save for %s: %w", userName, err)
	}

	return save, nil
}

func load(db *sql.DB, userID int) (string, error) {
	stmt, err := db.Prepare("select  state_json from saves where user_id = $1")
	if err != nil {
		return "", fmt.Errorf("prepare insert: %w", err)
	}

	defer stmt.Close()

	var save string

	err = stmt.QueryRow(userID).Scan(&save)

	if err != nil {
		return "", fmt.Errorf("error pulling user Id: %w", err)
	}

	return save, nil
}
