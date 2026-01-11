package handler

import (
	"encoding/json"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newThreadLikeRouter(threadRepo repository.ThreadRepository, likeRepo repository.ThreadLikeRepository, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := service.NewThreadLikeService(threadRepo, likeRepo, threadRepo)
	h := NewThreadLikeHandler(svc)

	r := gin.New()
	auth := r.Group("/api")
	auth.Use(testAuthMiddleware(userID))
	auth.POST("/threads/:id/like", h.Like)
	auth.DELETE("/threads/:id/like", h.Unlike)
	auth.GET("/threads/:id/like", h.Status)

	return r
}

func TestThreadLikeUnauthorized(t *testing.T) {
	r := newThreadLikeRouter(&fakeThreadRepo{}, &fakeThreadLikeRepo{}, 0)

	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestThreadLikeNotFound(t *testing.T) {
	r := newThreadLikeRouter(&fakeThreadRepo{findResult: nil}, &fakeThreadLikeRepo{}, 1)

	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestThreadLikeConflict(t *testing.T) {
	r := newThreadLikeRouter(
		&fakeThreadRepo{findResult: &models.Thread{Model: gorm.Model{ID: 1}, UserID: 1}},
		&fakeThreadLikeRepo{createErr: repository.ErrAlreadyLiked},
		1,
	)

	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusConflict, w.Code, w.Body.String())
	}
}

func TestThreadLikeOK(t *testing.T) {
	r := newThreadLikeRouter(
		&fakeThreadRepo{findResult: &models.Thread{Model: gorm.Model{ID: 1}, UserID: 1}},
		&fakeThreadLikeRepo{},
		1,
	)

	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	var resp struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.Message == "" {
		t.Fatalf("expected message, got empty")
	}
}

func TestThreadUnlikeConflict(t *testing.T) {
	r := newThreadLikeRouter(
		&fakeThreadRepo{findResult: &models.Thread{Model: gorm.Model{ID: 1}, UserID: 1}},
		&fakeThreadLikeRepo{deleteErr: repository.ErrLikeNotFound},
		1,
	)

	req := httptest.NewRequest(http.MethodDelete, "/api/threads/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusConflict, w.Code, w.Body.String())
	}
}

func TestThreadLikeStatusUnauthorized(t *testing.T) {
	r := newThreadLikeRouter(&fakeThreadRepo{}, &fakeThreadLikeRepo{}, 0)

	req := httptest.NewRequest(http.MethodGet, "/api/threads/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestThreadLikeStatusNotFound(t *testing.T) {
	r := newThreadLikeRouter(&fakeThreadRepo{findResult: nil}, &fakeThreadLikeRepo{}, 1)

	req := httptest.NewRequest(http.MethodGet, "/api/threads/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestThreadLikeStatusOK(t *testing.T) {
	r := newThreadLikeRouter(
		&fakeThreadRepo{findResult: &models.Thread{Model: gorm.Model{ID: 1}, UserID: 1}},
		&fakeThreadLikeRepo{exists: true},
		1,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/threads/1/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		Liked bool `json:"liked"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !resp.Liked {
		t.Fatalf("expected liked true, got false")
	}
}
