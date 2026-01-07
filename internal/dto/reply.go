package dto

import "time"

type CreateReplyReq struct {
	Content string `json:"content" binding:"required,min=1,max=4000"`
}

type ReplyResp struct {
	ID        uint      `json:"id"`
	ThreadID  uint      `json:"thread_id"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
