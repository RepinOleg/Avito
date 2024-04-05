package model

type BannerBody struct {
	// Идентификатор фичи
	FeatureID int64 `json:"feature_id,omitempty"`

	// Идентификаторы тэгов
	TagIDs []int64 `json:"tag_ids,omitempty"`

	// Содержимое баннера
	Content BannerContent `json:"content,omitempty"`

	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`
}

type BannerContent struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
	URL   string `json:"url,omitempty"`
}
