package api

type LoginRequest struct {
	AdminID  string `json:"adminID"`
	Password string `json:"password"`
}
