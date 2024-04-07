package memorycache

import (
	"errors"
	"github.com/RepinOleg/Banner_service/internal/model"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	banners           map[int64]model.BannerBody
}

func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	items := make(map[int64]model.BannerBody)

	cache := Cache{
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
		banners:           items,
	}

	if cleanupInterval > 0 {
		cache.StartGC() // данный метод рассматривается ниже
	}

	return &cache
}

func (c *Cache) Set(id int64, item model.BannerBody, duration time.Duration) {
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

func (c *Cache) GetBanner(tagID, featureID int64) ([]model.BannerContent, error) {
	c.RLock()

	defer c.RUnlock()
	var banners []model.BannerContent
	for _, banner := range c.banners {
		if banner.FeatureID == featureID {
			for _, tag := range banner.TagIDs {
				if tag == tagID {
					if banner.Expiration < 0 || time.Now().UnixNano() < banner.Expiration {
						banners = append(banners, banner.Content)
					}
				}
			}
		}
	}

	return banners, nil
}

func (c *Cache) Delete(id int64) error {

	c.Lock()

	defer c.Unlock()

	if _, ok := c.banners[id]; !ok {
		return errors.New("key not found")
	}

	delete(c.banners, id)

	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
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
func (c *Cache) expiredKeys() (ids []int64) {
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
func (c *Cache) clearItems(id []int64) {
	c.Lock()
	defer c.Unlock()

	for _, k := range id {
		delete(c.banners, k)
	}
}