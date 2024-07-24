package cache

import (
	"realty/db"
	"realty/models"
	"sync"
)

type UserCache struct {
	CurrentUser models.User
	OldUser     models.User
	ToCreate    bool
	ToUpdate    bool
	ToDelete    bool
	Deleted     bool
	mu          sync.RWMutex
}

func (user *UserCache) Save() error {
	user.mu.Lock()
	defer user.mu.Unlock()
	if user.Deleted {
		return nil
	}
	if user.ToDelete {
		err := db.DeleteAdv(user.CurrentUser.Id)
		if err != nil {
			return err
		}
		user.Deleted = true
		user.ToDelete = false
		user.ToCreate = false
		user.ToUpdate = false
	}
	if user.ToCreate {
		err := db.CreateUser(&user.CurrentUser)
		if err != nil {
			return err
		}
		user.OldUser = user.CurrentUser
		user.ToCreate = false
		user.ToUpdate = false
	}
	if user.ToUpdate {
		err := db.UpdateUserChanges(&user.OldUser, &user.CurrentUser)
		if err != nil {
			return err
		}
		user.OldUser = user.CurrentUser
		user.ToUpdate = false
	}
	return nil
}
