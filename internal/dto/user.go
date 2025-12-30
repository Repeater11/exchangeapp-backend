package dto

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResp struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
}
