package responses

import (
	"github.com/RepinOleg/Banner_service/internal/models"
	"time"
)

type Response200 struct {
	// Идентификатор баннера
	BannerId int32 `json:"banner_id,omitempty"`
	// Идентификаторы тэгов
	TagIds []int32 `json:"tag_ids,omitempty"`
	// Идентификатор фичи
	FeatureId int32 `json:"feature_id,omitempty"`
	// Содержимое баннера
	Content models.BannerContent `json:"content,omitempty"`
	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`
	// Дата создания баннера
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Дата обновления баннера
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
