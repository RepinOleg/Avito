package response

import (
	"github.com/RepinOleg/Banner_service/internal/model"
	"time"
)

type ModelResponse200 struct {
	// Идентификатор баннера
	BannerID int64 `json:"banner_id,omitempty"`
	// Идентификаторы тэгов
	TagIDs []int64 `json:"tag_ids,omitempty"`
	// Идентификатор фичи
	FeatureID int64 `json:"feature_id,omitempty"`
	// Содержимое баннера
	Content model.BannerContent `json:"content,omitempty"`
	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`
	// Дата создания баннера
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Дата обновления баннера
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
