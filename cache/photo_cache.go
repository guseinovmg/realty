package cache

import (
	"realty/db"
	"realty/models"
	"sync"
)

type PhotoCache struct {
	Photo    models.Photo
	ToCreate bool
	ToDelete bool
	Deleted  bool
	mu       sync.RWMutex
}

func (photo *PhotoCache) Save() error {
	photo.mu.Lock()
	defer photo.mu.Unlock()
	if photo.Deleted {
		return nil
	}
	if photo.ToDelete {
		err := db.DeletePhoto(photo.Photo.Id)
		if err != nil {
			return err
		}
		photo.Deleted = true
		photo.ToDelete = false
		photo.ToCreate = false
	}
	if photo.ToCreate {
		err := db.CreatePhoto(photo.Photo)
		if err != nil {
			return err
		}
		photo.ToCreate = false
	}
	return nil
}
