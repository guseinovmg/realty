package db

import (
	"database/sql"
	"log"
	"realty/config"
	"realty/models"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Initialize() {
	db, err := sql.Open("sqlite", config.GetDbPath())
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query(`SELECT * FROM users;`)
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	for rows.Next() {
		var (
			id int64
		)
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		log.Printf("id %d", id)
	}
	result, err := db.Exec(`INSERT INTO users (id) VALUES (88)`)
	if err != nil {
		log.Fatal(err)
	}
	li, _ := result.LastInsertId()
	ra, _ := result.RowsAffected()
	log.Println(li, ra)
	DB = db
}

func ReadDb() ([]models.User, []models.Adv, error) {
	return nil, nil, nil
}

func InsertAdv(adv *models.Adv) error {
	return nil
}

func GetAdv(id int64) (*models.Adv, error) {
	return nil, nil
}

func UpdateAdv(adv *models.Adv) error {
	return nil
}

func DeleteAdv(id int64) error {
	return nil
}

func InsertUser(user *models.User) error {
	return nil
}

func GetUser(id int64) (*models.User, error) {
	return nil, nil
}

func UpdateUser(user *models.User) error {
	return nil
}

func DeleteUser(id int64) error {
	return nil
}
