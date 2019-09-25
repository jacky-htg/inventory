package request

//LoginRequest : format json request for login
type LoginRequest struct {
	Username string `json:"username"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}
