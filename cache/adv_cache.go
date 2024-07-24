package cache

import (
	"realty/db"
	"realty/models"
	"strconv"
	"sync"
)

type AdvCache struct {
	CurrentAdv models.Adv
	OldAdv     models.Adv
	Photos     []*PhotoCache
	Watches    *WatchesCache
	ToCreate   bool
	ToUpdate   bool
	ToDelete   bool
	Deleted    bool
	mu         sync.RWMutex
	photoMu    sync.RWMutex
}

func (adv *AdvCache) Save() error {
	adv.mu.Lock()
	defer adv.mu.Unlock()
	if adv.Deleted {
		return nil
	}
	if adv.ToDelete {
		err := db.DeleteAdv(adv.CurrentAdv.Id)
		if err != nil {
			return err
		}
		adv.Deleted = true
		adv.ToDelete = false
		adv.ToCreate = false
		adv.ToUpdate = false
	}
	if adv.ToCreate {
		err := db.CreateAdv(&adv.CurrentAdv)
		if err != nil {
			return err
		}
		adv.OldAdv = adv.CurrentAdv
		adv.ToCreate = false
		adv.ToUpdate = false
	}
	if adv.ToUpdate {
		err := db.UpdateAdvChanges(&adv.OldAdv, &adv.CurrentAdv)
		if err != nil {
			return err
		}
		adv.OldAdv = adv.CurrentAdv
		adv.ToUpdate = false
	}
	return nil
}

func (adv *AdvCache) GetPhotosFilenames() []string {
	result := make([]string, 0, len(adv.Photos))
	adv.photoMu.RLock()
	defer adv.photoMu.RUnlock()
	for _, v := range adv.Photos {
		if v.Deleted || v.ToDelete {
			continue
		}
		ext := ""
		switch v.Photo.Ext {
		case 1:
			ext = ".jpg"
		case 2:
			ext = ".png"
		case 3:
			ext = ".gif"
		}
		name := strconv.FormatInt(v.Photo.Id, 10) + ext
		result = append(result, name)
	}
	return result
}
