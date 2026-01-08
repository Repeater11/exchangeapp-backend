package dto

import "time"

type CreateThreadReq struct {
	Title   string `json:"title" binding:"required,min=1,max=200"`
	Content string `json:"content" binding:"required,min=1,max=10000"`
}

type ThreadSummaryResp struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ThreadDetailResp struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ThreadListResp struct {
	Items []ThreadSummaryResp `json:"items"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}
