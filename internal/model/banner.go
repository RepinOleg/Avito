package model

import "time"

type BannerBody struct {
	// Идентификаторы тэгов
	TagIDs []int64 `json:"tag_ids,omitempty"`

	// Идентификатор фичи
	FeatureID int64 `json:"feature_id,omitempty"`

	// Содержимое баннера
	Content BannerContent `json:"content,omitempty"`

	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`

	Expiration int64
}

type BannerContent struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
}
