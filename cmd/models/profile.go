package models

type ProfilePictureRequest struct {
	Email string `json:"email"`
}

type ProfilePictureResponse struct {
	Email     string `json:"email"`
	Image     string `json:"image"`
	IsTecheer bool   `json:"isTecheer"`
}
