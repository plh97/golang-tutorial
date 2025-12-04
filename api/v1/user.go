package v1

import "go-nunu/internal/model"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Name  string `json:"name" example:"alan"`
	Email string `json:"email" example:"1234@gmail.com"`
	Image string `json:"image"`
}
type GetProfileResponseData struct {
	UserId string `json:"userId"`
	Name   string `json:"name" example:"alan"`
	Email  string `json:"email" example:"alan"`
	Image  string `json:"image"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type GetUserListResponseData struct {
	List []model.User `json:"list"`
}
type GetUserListResponse struct {
	Response
	Data GetUserListResponseData
}

type UpdateUserRequest struct {
	UserId  string `json:"user_id" binding:"required"`
	Name    string `json:"name" example:"alan"`
	Email   string `json:"email"`
	Image   string `json:"image"`
	RoleIds []uint `json:"role_ids"`
}
