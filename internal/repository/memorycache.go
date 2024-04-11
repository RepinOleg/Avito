package repository

import (
	"sync"
	"time"

	"github.com/RepinOleg/Banner_service/internal/model"
	"github.com/RepinOleg/Banner_service/internal/response"
)

type MemoryCache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	banners           map[int64]model.BannerBody
}

func New(defaultExpiration, cleanupInterval time.Duration) *MemoryCache {
	items := make(map[int64]model.BannerBody)

	cache := MemoryCache{
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		banners:           items,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

func (c *MemoryCache) AddBanner(id int64, item model.BannerBody, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()
	item.CreatedAt = time.Now()
	item.Expiration = expiration
	c.banners[id] = item
}

func (c *MemoryCache) GetBanner(tagID, featureID int64) (*model.BannerContent, error) {
	c.RLock()

	defer c.RUnlock()
	for _, banner := range c.banners {
		if banner.FeatureID != featureID {
			continue
		}

		for _, tag := range banner.TagIDs {
			if tag != tagID {
				continue
			}

			if banner.Expiration < 0 || time.Now().UnixNano() < banner.Expiration {
				return &banner.Content, nil
			}
		}
	}

	return nil, &response.NotFoundError{Message: "banner not found"}
}

func (c *MemoryCache) StartGC() {
	go c.GC()
}

func (c *MemoryCache) GC() {
	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.banners == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)

		}
	}
}

// expiredKeys возвращает список "просроченных" ключей
func (c *MemoryCache) expiredKeys() (ids []int64) {
	c.RLock()
	defer c.RUnlock()

	for k, i := range c.banners {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			ids = append(ids, k)
		}
	}
	return
}

// clearItems удаляет ключи из переданного списка
func (c *MemoryCache) clearItems(id []int64) {
	c.Lock()
	defer c.Unlock()

	for _, k := range id {
		delete(c.banners, k)
	}
}
