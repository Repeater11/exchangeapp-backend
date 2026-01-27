package handler

import (
	"encoding/json"
	"errors"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/service"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newReplyRouter(replyRepo repository.ReplyRepository, threadRepo repository.ThreadRepository, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := service.NewReplyService(replyRepo, threadRepo)
	h := NewReplyHandler(svc)

	r := gin.New()
	r.GET("/threads/:id/replies", h.ListByThreadID)

	auth := r.Group("/api")
	auth.Use(testAuthMiddleware(userID))
	auth.GET("/me/replies", h.ListMine)
	auth.POST("/threads/:id/replies", h.Create)
	auth.PUT("/replies/:id", h.Update)
	auth.DELETE("/replies/:id", h.Delete)

	return r
}

func TestReplyListNotFound(t *testing.T) {
	replyRepo := &fakeReplyRepo{}
	threadRepo := &fakeThreadRepo{findResult: nil}
	r := newReplyRouter(replyRepo, threadRepo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads/1/replies", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestReplyListBadParam(t *testing.T) {
	replyRepo := &fakeReplyRepo{}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads/abc/replies", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestReplyListOK(t *testing.T) {
	replyRepo := &fakeReplyRepo{
		listResult: []models.Reply{
			{Model: gorm.Model{ID: 1}, ThreadID: 1, UserID: 2, Content: "c"},
		},
		countResult: 1,
	}
	threadRepo := &fakeThreadRepo{
		findResult: &models.Thread{ID: 1},
	}
	r := newReplyRouter(replyRepo, threadRepo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads/1/replies?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp dto.ReplyListResp
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestReplyListCursorInvalid(t *testing.T) {
	replyRepo := &fakeReplyRepo{}
	threadRepo := &fakeThreadRepo{findResult: &models.Thread{ID: 1}}
	r := newReplyRouter(replyRepo, threadRepo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads/1/replies?cursor=bad", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestReplyListCursorOK(t *testing.T) {
	ts := time.Unix(0, 123)
	replyRepo := &fakeReplyRepo{
		listAfterResult: []models.Reply{
			{Model: gorm.Model{ID: 7, CreatedAt: ts}, ThreadID: 1, UserID: 2, Content: "c"},
		},
	}
	threadRepo := &fakeThreadRepo{findResult: &models.Thread{ID: 1}}
	r := newReplyRouter(replyRepo, threadRepo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads/1/replies?cursor=1_1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp dto.ReplyListResp
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	wantCursor := strconv.FormatInt(ts.UnixNano(), 10) + "_7"
	if resp.NextCursor != wantCursor {
		t.Fatalf("expected next_cursor %s, got %s", wantCursor, resp.NextCursor)
	}
}

func TestReplyCreateUnauthorized(t *testing.T) {
	replyRepo := &fakeReplyRepo{}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 0)

	body := `{"content":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/replies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestReplyCreateBadBody(t *testing.T) {
	replyRepo := &fakeReplyRepo{}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	body := `{"content":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/replies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestReplyCreateThreadNotFound(t *testing.T) {
	replyRepo := &fakeReplyRepo{}
	threadRepo := &fakeThreadRepo{findResult: nil}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	body := `{"content":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/replies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestReplyCreateRepoError(t *testing.T) {
	replyRepo := &fakeReplyRepo{createErr: errors.New("boom")}
	threadRepo := &fakeThreadRepo{findResult: &models.Thread{ID: 1}}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	body := `{"content":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/replies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

func TestReplyCreateOK(t *testing.T) {
	replyRepo := &fakeReplyRepo{}
	threadRepo := &fakeThreadRepo{
		findResult: &models.Thread{ID: 1},
	}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	body := `{"content":"hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads/1/replies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}

	var resp dto.ReplyResp
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.UserID != 1 || resp.ThreadID != 1 || resp.Content != "hello" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestReplyListMineUnauthorized(t *testing.T) {
	replyRepo := &fakeReplyRepo{}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 0)

	req := httptest.NewRequest(http.MethodGet, "/api/me/replies", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestReplyListMineOK(t *testing.T) {
	replyRepo := &fakeReplyRepo{
		listResult: []models.Reply{
			{Model: gorm.Model{ID: 1}, ThreadID: 1, UserID: 1, Content: "c"},
		},
		countResult: 1,
	}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	req := httptest.NewRequest(http.MethodGet, "/api/me/replies?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp dto.ReplyListResp
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestReplyListMineCursorOK(t *testing.T) {
	ts := time.Unix(0, 456)
	replyRepo := &fakeReplyRepo{
		listByUserAfterRes: []models.Reply{
			{Model: gorm.Model{ID: 9, CreatedAt: ts}, ThreadID: 1, UserID: 1, Content: "c"},
		},
	}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	req := httptest.NewRequest(http.MethodGet, "/api/me/replies?cursor=1_1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp dto.ReplyListResp
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	wantCursor := strconv.FormatInt(ts.UnixNano(), 10) + "_9"
	if resp.NextCursor != wantCursor {
		t.Fatalf("expected next_cursor %s, got %s", wantCursor, resp.NextCursor)
	}
}

func TestReplyUpdateNotFound(t *testing.T) {
	replyRepo := &fakeReplyRepo{findResult: nil}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	body := `{"content":"hello"}`
	req := httptest.NewRequest(http.MethodPut, "/api/replies/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestReplyUpdateForbidden(t *testing.T) {
	replyRepo := &fakeReplyRepo{
		findResult: &models.Reply{Model: gorm.Model{ID: 1}, UserID: 2, ThreadID: 1},
	}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	body := `{"content":"hello"}`
	req := httptest.NewRequest(http.MethodPut, "/api/replies/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusForbidden, w.Code, w.Body.String())
	}
}

func TestReplyDeleteNotFound(t *testing.T) {
	replyRepo := &fakeReplyRepo{findResult: nil}
	threadRepo := &fakeThreadRepo{}
	r := newReplyRouter(replyRepo, threadRepo, 1)

	req := httptest.NewRequest(http.MethodDelete, "/api/replies/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}
