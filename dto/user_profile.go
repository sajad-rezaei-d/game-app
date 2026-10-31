package dto


type ProfileRequest struct {
	UserID uint	`json:"id"`
}

type ProfileResponse struct {
	Name string `json:"name"`
}