package dto

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type RegisterResp struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
}
