package response

import (
	"github.com/RepinOleg/Banner_service/internal/model"
	"time"
)

type Response200 struct {
	// Идентификатор баннера
	BannerID int32 `json:"banner_id,omitempty"`
	// Идентификаторы тэгов
	TagIDs []int32 `json:"tag_ids,omitempty"`
	// Идентификатор фичи
	FeatureID int32 `json:"feature_id,omitempty"`
	// Содержимое баннера
	Content model.BannerContent `json:"content,omitempty"`
	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`
	// Дата создания баннера
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Дата обновления баннера
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
