package models

type BannerBody struct {
	// Идентификаторы тэгов
	TagIds []int32 `json:"tag_ids,omitempty"`
	// Идентификатор фичи
	FeatureId int32 `json:"feature_id,omitempty"`
	// Содержимое баннера
	Content ModelMap `json:"content,omitempty"`
	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`
}
