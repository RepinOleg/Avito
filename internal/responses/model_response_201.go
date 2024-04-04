package responses

type Response201 struct {
	// Идентификатор созданного баннера
	BannerId int32 `json:"banner_id,omitempty"`
}
