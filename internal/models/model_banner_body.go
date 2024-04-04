package models

type BannerBody struct {
	// Идентификатор фичи
	FeatureId int32 `json:"feature_id,omitempty"`

	// Идентификаторы тэгов
	TagIds []int32 `json:"tag_ids,omitempty"`

	// Содержимое баннера
	Content BannerContent `json:"content,omitempty"`

	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`
}
