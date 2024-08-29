package db

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"realty/config"
	"realty/models"
	"strings"

	_ "modernc.org/sqlite"
)

var dbUsers, dbAdvs, dbPhotos, dbWatches *sql.DB

func Initialize() {
	if db, err := sql.Open("sqlite", config.GetDbUsersPath()); err != nil {
		log.Fatal(err)
	} else {
		db.SetMaxOpenConns(10)
		dbUsers = db
	}
	if db, err := sql.Open("sqlite", config.GetDbAdvsPath()); err != nil {
		log.Fatal(err)
	} else {
		db.SetMaxOpenConns(1)
		dbAdvs = db
	}
	if db, err := sql.Open("sqlite", config.GetDbPhotosPath()); err != nil {
		log.Fatal(err)
	} else {
		db.SetMaxOpenConns(1)
		dbPhotos = db
	}
	if db, err := sql.Open("sqlite", config.GetDbWatchesPath()); err != nil {
		log.Fatal(err)
	} else {
		db.SetMaxOpenConns(1)
		dbWatches = db
	}
	if config.GetDataDir() == ":memory:" {
		if err := CreateInMemoryDB(); err != nil {
			log.Fatal(err)
		}
	}
}

func CreateInMemoryDB() error {
	if _, err := dbUsers.Exec(`create table users
(
    id             INTEGER
        primary key,
    email          TEXT      not null
        unique,
    name           TEXT      not null,
    password_hash  BLOB      not null,
    session_secret BLOB      not null,
    invite_id      TEXT,
    balance        REAL      not null,
    trusted        INTEGER   not null,
    enabled        INTEGER   not null,
    description    TEXT
) without ROWID, strict;`); err != nil {
		return errors.Join(err, errors.New("db.CreateInMemoryDB() 1"))
	}

	if _, err := dbUsers.Exec(`create table invites
(
    id   TEXT primary key,
	name TEXT
) without ROWID, strict;`); err != nil {
		return errors.Join(err, errors.New("db.CreateInMemoryDB() 2"))
	}

	if _, err := dbAdvs.Exec(`
		    CREATE TABLE advs (
		        id INTEGER PRIMARY KEY,
		        user_id INTEGER NOT NULL,
		        updated INTEGER NOT NULL,
		        approved INTEGER NOT NULL,
		        lang INTEGER NOT NULL,
		        origin_lang INTEGER NOT NULL,
		        translated_by INTEGER NOT NULL,
		        translated_to TEXT NOT NULL,
		        title TEXT NOT NULL,
		        description TEXT NOT NULL,
		        price INTEGER NOT NULL,
		        currency TEXT NOT NULL,
		        country TEXT NOT NULL,
		        city TEXT NOT NULL,
		        address TEXT NOT NULL,
		        latitude REAL NOT NULL,
		        longitude REAL NOT NULL,
		        paid_adv INTEGER NOT NULL,
		        se_visible INTEGER NOT NULL,
		        user_comment TEXT NOT NULL,
		        admin_comment TEXT NOT NULL
		    ) without ROWID, strict;
		`); err != nil {
		return errors.Join(err, errors.New("db.CreateInMemoryDB() 3"))
	}

	if _, err := dbPhotos.Exec(`
		    CREATE TABLE photos (
		        id INTEGER PRIMARY KEY,
		        adv_id INTEGER NOT NULL,
		        ext INTEGER NOT NULL
		    ) without ROWID, strict;
		`); err != nil {
		return errors.Join(err, errors.New("db.CreateInMemoryDB() 4"))
	}

	if _, err := dbWatches.Exec(`
    CREATE TABLE watches (
        adv_id INTEGER PRIMARY KEY,
        count INTEGER NOT NULL
    ) without ROWID, strict;
`); err != nil {
		return errors.Join(err, errors.New("db.CreateInMemoryDB() 5"))
	}
	return nil
}

func ReadDb() (users []*models.User, advs []*models.Adv, photos []*models.Photo, watches []*models.Watches, err error) {
	if users, err = GetUsers(); err != nil {
		return
	}
	if advs, err = GetAdvs(); err != nil {
		return
	}
	if photos, err = GetPhotos(); err != nil {
		return
	}
	if watches, err = GetWatches(); err != nil {
		return
	}
	return
}

func CreateAdv(adv *models.Adv) error {
	query := `
		INSERT INTO advs (
			id, user_id, updated, approved, lang, origin_lang, title,
			description, price, currency, country, city, address, latitude,
			longitude, paid_adv, se_visible, user_comment,
			admin_comment, translated_to
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := dbAdvs.Exec(query,
		adv.Id, adv.UserId, adv.Updated, adv.Approved, adv.Lang,
		adv.OriginLang, adv.Title, adv.Description, adv.Price, adv.Currency,
		adv.Country, adv.City, adv.Address, adv.Latitude, adv.Longitude,
		adv.PaidAdv, adv.SeVisible, adv.UserComment,
		adv.AdminComment, adv.TranslatedTo,
	)
	if err != nil {
		return errors.Join(err, errors.New("db.CreateAdv()"))
	}

	return nil
}

func GetAdv(id int64) (*models.Adv, error) {
	adv := &models.Adv{}
	query := "SELECT * FROM advs WHERE id = ?"
	err := dbAdvs.QueryRow(query, id).Scan(
		&adv.Id, &adv.UserId, &adv.Updated, &adv.Approved,
		&adv.Lang, &adv.OriginLang, &adv.Title, &adv.Description, &adv.Price,
		&adv.Currency, &adv.Country, &adv.City, &adv.Address, &adv.Latitude,
		&adv.Longitude, &adv.PaidAdv, &adv.SeVisible,
		&adv.UserComment, &adv.AdminComment, &adv.TranslatedTo,
	)
	if err != nil {
		return nil, errors.Join(err, errors.New("db.GetAdv()"))
	}

	return adv, nil
}

func GetAdvs() ([]*models.Adv, error) {
	rows, err := dbAdvs.Query("SELECT * FROM advs ORDER BY id")
	if err != nil {
		return nil, errors.Join(err, errors.New("db.GetAdvs()"))
	}
	defer rows.Close() //todo нужо ли закрывать соединение?

	var advs = make([]*models.Adv, 0, 1000)

	for rows.Next() {
		adv := &models.Adv{}
		err := rows.Scan(
			&adv.Id, &adv.UserId, &adv.Updated, &adv.Approved,
			&adv.Lang, &adv.OriginLang, &adv.Title, &adv.Description, &adv.Price,
			&adv.Currency, &adv.Country, &adv.City, &adv.Address, &adv.Latitude,
			&adv.Longitude, &adv.PaidAdv, &adv.SeVisible,
			&adv.UserComment, &adv.AdminComment, &adv.TranslatedTo,
		)
		if err != nil {
			return nil, errors.Join(err, errors.New("db.GetAdvs()"))
		}
		advs = append(advs, adv)
	}

	return advs, nil
}

func UpdateAdv(adv *models.Adv) error {
	query := `
		UPDATE advs SET
			user_id = ?,
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
			paid_adv = ?,
			se_visible = ?,
			user_comment = ?,
			admin_comment = ?,
			translated_to = ?
		WHERE id = ?
	`
	_, err := dbUsers.Exec(query,
		adv.UserId, adv.Updated, adv.Approved, adv.Lang,
		adv.OriginLang, adv.Title, adv.Description, adv.Price, adv.Currency,
		adv.Country, adv.City, adv.Address, adv.Latitude, adv.Longitude,
		adv.PaidAdv, adv.SeVisible, adv.UserComment,
		adv.AdminComment, adv.TranslatedTo, adv.Id,
	)
	if err != nil {
		return errors.Join(err, errors.New("db.UpdateAdv()"))
	}
	return nil
}

func UpdateAdvChanges(oldAdv, newAdv *models.Adv) error {
	args := make([]interface{}, 0, 16)
	setClauses := make([]string, 0, 16)

	if oldAdv.UserId != newAdv.UserId {
		setClauses = append(setClauses, "user_id = ?")
		args = append(args, newAdv.UserId)
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
	if oldAdv.PaidAdv != newAdv.PaidAdv {
		setClauses = append(setClauses, "paid_adv = ?")
		args = append(args, newAdv.PaidAdv)
	}
	if oldAdv.SeVisible != newAdv.SeVisible {
		setClauses = append(setClauses, "se_visible = ?")
		args = append(args, newAdv.SeVisible)
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

	if len(setClauses) == 0 {
		return nil
	}

	query := "UPDATE advs SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, oldAdv.Id)

	_, err := dbAdvs.Exec(query, args...)

	if err != nil {
		return errors.Join(err, errors.New("db.UpdateAdvChanges()"))
	}
	return nil
}

func DeleteAdv(id int64) error {
	query := "DELETE FROM advs WHERE id = ?"
	_, err := dbAdvs.Exec(query, id)
	if err != nil {
		return errors.Join(err, errors.New("db.DeleteAdv()"))
	}
	return nil
}

func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (
			id, email, name, password_hash, session_secret, invite_id, trusted,
			enabled, balance, description
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`
	_, err := dbUsers.Exec(query,
		user.Id, user.Email, user.Name, user.PasswordHash, user.SessionSecret[:],
		user.InviteId, user.Trusted, user.Enabled, user.Balance,
		user.Description,
	)
	if err != nil {
		return errors.Join(err, errors.New("db.CreateUser()"))
	}

	return nil
}

func GetUser(id int64) (*models.User, error) {
	user := &models.User{}
	query := "SELECT * FROM users WHERE id = ?"
	err := dbUsers.QueryRow(query, id).Scan(
		&user.Id, &user.Email, &user.Name, &user.PasswordHash,
		&user.SessionSecret, &user.InviteId, &user.Trusted, &user.Enabled,
		&user.Balance, &user.Description,
	)
	if err != nil {
		return nil, errors.Join(err, errors.New("db.GetUser()"))
	}

	return user, nil
}

func UpdateUser(user *models.User) error {
	query := `
		UPDATE users SET
			email = ?,
			name = ?,
			password_hash = ?,
			session_secret = ?,
			invite_id = ?,
			trusted = ?,
			enabled = ?,
			balance = ?,
			description = ?
		WHERE id = ?
	`
	_, err := dbUsers.Exec(query,
		user.Email, user.Name, user.PasswordHash, user.SessionSecret[:],
		user.InviteId, user.Trusted, user.Enabled, user.Balance,
		user.Description, user.Id,
	)

	if err != nil {
		return errors.Join(err, errors.New("db.UpdateUser()"))
	}
	return nil
}

func UpdateUserChanges(oldUser, newUser *models.User) error {

	args := make([]interface{}, 0, 9)
	setClauses := make([]string, 0, 9)

	if oldUser.Email != newUser.Email {
		setClauses = append(setClauses, "email = ?")
		args = append(args, newUser.Email)
	}
	if oldUser.Name != newUser.Name {
		setClauses = append(setClauses, "name = ?")
		args = append(args, newUser.Name)
	}
	if !bytes.Equal(oldUser.PasswordHash, newUser.PasswordHash) {
		setClauses = append(setClauses, "password_hash = ?")
		args = append(args, newUser.PasswordHash)
	}
	if !bytes.Equal(oldUser.SessionSecret[:], newUser.SessionSecret[:]) {
		setClauses = append(setClauses, "session_secret = ?")
		args = append(args, newUser.SessionSecret[:])
	}
	if oldUser.InviteId != newUser.InviteId {
		setClauses = append(setClauses, "invite_id = ?")
		args = append(args, newUser.InviteId)
	}
	if oldUser.Trusted != newUser.Trusted {
		setClauses = append(setClauses, "trusted = ?")
		args = append(args, newUser.Trusted)
	}
	if oldUser.Enabled != newUser.Enabled {
		setClauses = append(setClauses, "enabled = ?")
		args = append(args, newUser.Enabled)
	}
	if oldUser.Balance != newUser.Balance {
		setClauses = append(setClauses, "balance = ?")
		args = append(args, newUser.Balance)
	}
	if oldUser.Description != newUser.Description {
		setClauses = append(setClauses, "description = ?")
		args = append(args, newUser.Description)
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := "UPDATE users SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, oldUser.Id)

	_, err := dbUsers.Exec(query, args...)
	if err != nil {
		return errors.Join(err, errors.New("db.UpdateUserChanges()"))
	}
	return nil
}

func GetUsers() ([]*models.User, error) {
	rows, err := dbUsers.Query("SELECT * FROM users ORDER BY id")
	if err != nil {
		return nil, errors.Join(err, errors.New("db.GetUsers()"))
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.Id, &user.Email, &user.Name, &user.PasswordHash,
			&user.SessionSecret, &user.InviteId, &user.Trusted, &user.Enabled,
			&user.Balance, &user.Description,
		)
		if err != nil {
			return nil, errors.Join(err, errors.New("db.GetUsers()"))
		}
		users = append(users, user)
	}

	return users, nil
}

func DeleteUser(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := dbUsers.Exec(query, id)
	if err != nil {
		return errors.Join(err, errors.New("db.DeleteUser()"))
	}
	return nil
}

func CreatePhoto(photo models.Photo) error {
	query := `
		INSERT INTO photos (
			id, adv_id, ext
		) VALUES (
			?, ?, ?
		)
	`
	_, err := dbPhotos.Exec(query,
		photo.Id, photo.AdvId, photo.Ext,
	)
	if err != nil {
		return errors.Join(err, errors.New("db.CreatePhoto()"))
	}

	return nil
}

func GetPhotos() ([]*models.Photo, error) {
	rows, err := dbPhotos.Query("SELECT id, adv_id, ext FROM photos ORDER BY id")
	if err != nil {
		return nil, errors.Join(err, errors.New("db.GetPhotos()"))
	}
	defer rows.Close()

	var photos []*models.Photo

	for rows.Next() {
		photo := &models.Photo{}
		err := rows.Scan(
			&photo.Id, &photo.AdvId, &photo.Ext,
		)
		if err != nil {
			return nil, errors.Join(err, errors.New("db.GetPhotos()"))
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func DeletePhoto(id int64) error {
	query := "DELETE FROM photos WHERE id = ?"
	_, err := dbPhotos.Exec(query, id)
	if err != nil {
		return errors.Join(err, errors.New("db.DeletePhoto()"))
	}
	return nil
}

func CreateWatches(watches models.Watches) error {
	query := `
		INSERT INTO watches (
			adv_id, count
		) VALUES (
			?, ?
		)
	`
	_, err := dbWatches.Exec(query,
		watches.AdvId, watches.Count,
	)
	if err != nil {
		return errors.Join(err, errors.New("db.CreateWatches()"))
	}

	return nil
}

func GetWatches() ([]*models.Watches, error) {
	rows, err := dbWatches.Query("SELECT adv_id, count FROM watches ORDER BY adv_id")
	if err != nil {
		return nil, errors.Join(err, errors.New("db.GetWatches()"))
	}
	defer rows.Close()
	var watches []*models.Watches
	for rows.Next() {
		watch := &models.Watches{}
		err := rows.Scan(
			&watch.AdvId, &watch.Count,
		)
		if err != nil {
			return nil, errors.Join(err, errors.New("db.GetWatches()"))
		}
		watches = append(watches, watch)
	}
	return watches, nil
}

func UpdateWatches(watch *models.Watches) error {
	query := `
		UPDATE watches SET
			count = ?
		WHERE id = ?
	`
	_, err := dbWatches.Exec(query,
		watch.Count, watch.AdvId,
	)

	if err != nil {
		return errors.Join(err, errors.New("db.UpdateWatches()"))
	}
	return nil
}

func DeleteWatches(id int64) error {
	query := "DELETE FROM watches WHERE adv_id = ?"
	_, err := dbWatches.Exec(query, id)
	if err != nil {
		return errors.Join(err, errors.New("db.DeleteWatches()"))
	}
	return nil
}
