package db

import (
	"database/sql"
	"log"
	"realty/config"
	"realty/models"
	"strings"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Initialize() {
	db_, err := sql.Open("sqlite", config.GetDbPath())
	if err != nil {
		log.Fatal(err)
	}
	db = db_
}

func ReadDb() ([]models.User, []models.Adv, error) {
	return nil, nil, nil
}

func CreateAdv(adv *models.Adv) error {
	query := `
		INSERT INTO advs (
			id, user_id, created, updated, approved, lang, origin_lang, title,
			description, price, currency, country, city, address, latitude,
			longitude, watches, paid_adv, visible_for_search_engines, user_comment,
			admin_comment, translated_to, photos
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	result, err := db.Exec(query,
		adv.Id, adv.UserId, adv.Created, adv.Updated, adv.Approved, adv.Lang,
		adv.OriginLang, adv.Title, adv.Description, adv.Price, adv.Currency,
		adv.Country, adv.City, adv.Address, adv.Latitude, adv.Longitude,
		adv.Watches, adv.PaidAdv, adv.VisibleForSearchEngines, adv.UserComment,
		adv.AdminComment, adv.TranslatedTo, adv.Photos,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	adv.Id = id

	return nil
}

func GetAdv(id int64) (*models.Adv, error) {
	adv := &models.Adv{}
	query := "SELECT * FROM advs WHERE id = ?"
	err := db.QueryRow(query, id).Scan(
		&adv.Id, &adv.UserId, &adv.Created, &adv.Updated, &adv.Approved,
		&adv.Lang, &adv.OriginLang, &adv.Title, &adv.Description, &adv.Price,
		&adv.Currency, &adv.Country, &adv.City, &adv.Address, &adv.Latitude,
		&adv.Longitude, &adv.Watches, &adv.PaidAdv, &adv.VisibleForSearchEngines,
		&adv.UserComment, &adv.AdminComment, &adv.TranslatedTo, &adv.Photos,
	)
	if err != nil {
		return nil, err
	}

	return adv, nil
}

func GetAdvs() ([]*models.Adv, error) {
	rows, err := db.Query("SELECT * FROM advs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var advs = make([]*models.Adv, 0, 1000)

	for rows.Next() {
		adv := &models.Adv{}
		err := rows.Scan(
			&adv.Id, &adv.UserId, &adv.Created, &adv.Updated, &adv.Approved,
			&adv.Lang, &adv.OriginLang, &adv.Title, &adv.Description, &adv.Price,
			&adv.Currency, &adv.Country, &adv.City, &adv.Address, &adv.Latitude,
			&adv.Longitude, &adv.Watches, &adv.PaidAdv, &adv.VisibleForSearchEngines,
			&adv.UserComment, &adv.AdminComment, &adv.TranslatedTo, &adv.Photos,
		)
		if err != nil {
			return nil, err
		}
		advs = append(advs, adv)
	}

	return advs, nil
}

func UpdateAdv(adv *models.Adv) error {
	query := `
		UPDATE advs SET
			user_id = ?,
			created = ?,
			updated = ?,
			approved = ?,
			lang = ?,
			origin_lang = ?,
			title = ?,
			description = ?,
			price = ?,
			currency = ?,
			country = ?,
			city = ?,
			address = ?,
			latitude = ?,
			longitude = ?,
			watches = ?,
			paid_adv = ?,
			visible_for_search_engines = ?,
			user_comment = ?,
			admin_comment = ?,
			translated_to = ?,
			photos = ?
		WHERE id = ?
	`
	_, err := db.Exec(query,
		adv.UserId, adv.Created, adv.Updated, adv.Approved, adv.Lang,
		adv.OriginLang, adv.Title, adv.Description, adv.Price, adv.Currency,
		adv.Country, adv.City, adv.Address, adv.Latitude, adv.Longitude,
		adv.Watches, adv.PaidAdv, adv.VisibleForSearchEngines, adv.UserComment,
		adv.AdminComment, adv.TranslatedTo, adv.Photos, adv.Id,
	)

	return err
}

func UpdateAdvChanges(oldAdv, newAdv *models.Adv) error {
	args := make([]interface{}, 0, 16)
	setClauses := make([]string, 0, 16)

	if oldAdv.UserId != newAdv.UserId {
		setClauses = append(setClauses, "user_id = ?")
		args = append(args, newAdv.UserId)
	}
	if !oldAdv.Created.Equal(newAdv.Created) {
		setClauses = append(setClauses, "created = ?")
		args = append(args, newAdv.Created)
	}
	if !oldAdv.Updated.Equal(newAdv.Updated) {
		setClauses = append(setClauses, "updated = ?")
		args = append(args, newAdv.Updated)
	}
	if oldAdv.Approved != newAdv.Approved {
		setClauses = append(setClauses, "approved = ?")
		args = append(args, newAdv.Approved)
	}
	if oldAdv.Lang != newAdv.Lang {
		setClauses = append(setClauses, "lang = ?")
		args = append(args, newAdv.Lang)
	}
	if oldAdv.OriginLang != newAdv.OriginLang {
		setClauses = append(setClauses, "origin_lang = ?")
		args = append(args, newAdv.OriginLang)
	}
	if oldAdv.Title != newAdv.Title {
		setClauses = append(setClauses, "title = ?")
		args = append(args, newAdv.Title)
	}
	if oldAdv.Description != newAdv.Description {
		setClauses = append(setClauses, "description = ?")
		args = append(args, newAdv.Description)
	}
	if oldAdv.Price != newAdv.Price {
		setClauses = append(setClauses, "price = ?")
		args = append(args, newAdv.Price)
	}
	if oldAdv.Currency != newAdv.Currency {
		setClauses = append(setClauses, "currency = ?")
		args = append(args, newAdv.Currency)
	}
	if oldAdv.Country != newAdv.Country {
		setClauses = append(setClauses, "country = ?")
		args = append(args, newAdv.Country)
	}
	if oldAdv.City != newAdv.City {
		setClauses = append(setClauses, "city = ?")
		args = append(args, newAdv.City)
	}
	if oldAdv.Address != newAdv.Address {
		setClauses = append(setClauses, "address = ?")
		args = append(args, newAdv.Address)
	}
	if oldAdv.Latitude != newAdv.Latitude {
		setClauses = append(setClauses, "latitude = ?")
		args = append(args, newAdv.Latitude)
	}
	if oldAdv.Longitude != newAdv.Longitude {
		setClauses = append(setClauses, "longitude = ?")
		args = append(args, newAdv.Longitude)
	}
	if oldAdv.Watches != newAdv.Watches {
		setClauses = append(setClauses, "watches = ?")
		args = append(args, newAdv.Watches)
	}
	if oldAdv.PaidAdv != newAdv.PaidAdv {
		setClauses = append(setClauses, "paid_adv = ?")
		args = append(args, newAdv.PaidAdv)
	}
	if oldAdv.VisibleForSearchEngines != newAdv.VisibleForSearchEngines {
		setClauses = append(setClauses, "visible_for_search_engines = ?")
		args = append(args, newAdv.VisibleForSearchEngines)
	}
	if oldAdv.UserComment != newAdv.UserComment {
		setClauses = append(setClauses, "user_comment = ?")
		args = append(args, newAdv.UserComment)
	}
	if oldAdv.AdminComment != newAdv.AdminComment {
		setClauses = append(setClauses, "admin_comment = ?")
		args = append(args, newAdv.AdminComment)
	}
	if oldAdv.TranslatedTo != newAdv.TranslatedTo {
		setClauses = append(setClauses, "translated_to = ?")
		args = append(args, newAdv.TranslatedTo)
	}
	if oldAdv.Photos != newAdv.Photos {
		setClauses = append(setClauses, "photos = ?")
		args = append(args, newAdv.Photos)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := "UPDATE advs SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, oldAdv.Id)

	_, err := db.Exec(query, args...)
	return err
}

func DeleteAdv(id int64) error {
	query := "DELETE FROM advs WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}
