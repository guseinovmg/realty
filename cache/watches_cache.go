package cache

import (
	"realty/db"
	"realty/models"
	"sync"
)

type WatchesCache struct {
	Watches  models.Watches
	ToCreate bool
	ToUpdate bool
	ToDelete bool
	Deleted  bool
	mu       sync.RWMutex
}

func (watch *WatchesCache) Save() error {
	watch.mu.Lock()
	defer watch.mu.Unlock()
	if watch.Deleted {
		return nil
	}
	if watch.ToDelete {
		err := db.DeleteWatches(watch.Watches.AdvId)
		if err != nil {
			return err
		}
		watch.Deleted = true
		watch.ToDelete = false
		watch.ToCreate = false
		watch.ToUpdate = false
	}
	if watch.ToCreate {
		err := db.CreateWatches(watch.Watches)
		if err != nil {
			return err
		}
		watch.ToCreate = false
		watch.ToUpdate = false
	}
	if watch.ToUpdate {
		err := db.UpdateWatches(&watch.Watches)
		if err != nil {
			return err
		}
		watch.ToUpdate = false
	}
	return nil
}
