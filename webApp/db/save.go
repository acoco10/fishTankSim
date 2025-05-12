package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func Save(userName string, saveJson string) error {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	id, err := getUserID(db, userName)
	if err != nil {
		return err
	}

	exists, err := checkIfSaveExists(db, id)
	if err != nil {
		return err
	}

	if !exists {
		err = createSave(db, id, saveJson)
		if err != nil {
			return fmt.Errorf("error creating save: %w", err)
		}
		println("created save in saves table")
		return nil
	}

	err = updateSave(db, id, saveJson)
	if err != nil {
		log.Fatal(fmt.Errorf("error updating save: %w", err))
	}
	println("updated save in saves table")

	return nil
}

func updateSave(db *sql.DB, userID int, saveJson string) error {
	stmt, err := db.Prepare(
		"UPDATE saves SET state_json = $1 WHERE user_id = $2")
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(userID, saveJson)
	if err != nil {
		return fmt.Errorf("execute update: %w", err)
	}

	println(res)

	return nil
}

func createSave(db *sql.DB, userID int, saveJson string) error {
	stmt, err := db.Prepare("INSERT INTO saves (user_id, state_json) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(userID, saveJson)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	fmt.Println("save created")
	return nil
}

func checkIfSaveExists(db *sql.DB, userId int) (bool, error) {
	stmt, err := db.Prepare("SELECT EXISTS(SELECT 1 FROM saves WHERE user_id = $1)")
	if err != nil {
		return false, fmt.Errorf("prepare insert: %w", err)
	}

	defer stmt.Close()

	save := stmt.QueryRow(userId)
	var exists bool

	err = save.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking for save: %w", err)
	}

	return exists, nil
}

func getUserID(db *sql.DB, userName string) (int, error) {
	stmt, err := db.Prepare("select id from users where username = $1")
	if err != nil {
		return 0, fmt.Errorf("prepare insert: %w", err)
	}

	defer stmt.Close()
	var id int

	err = stmt.QueryRow(userName).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error pulling user Id %s: %w", userName, err)
	}
	if id == 0 {
		return 0, fmt.Errorf("no valid id found for user: %s", userName)
	}
	return id, nil
}
