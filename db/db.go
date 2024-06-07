package db

import (
	"realty/models"
)

func Initialize() {

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
