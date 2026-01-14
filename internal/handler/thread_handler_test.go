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

func newThreadRouter(repo repository.ThreadRepository, userID uint) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := service.NewThreadService(repo, &fakeThreadLikeRepo{}, repo)
	h := NewThreadHandler(svc)

	r := gin.New()
	r.GET("/threads", h.List)
	r.GET("/threads/:id", h.Detail)

	auth := r.Group("/api")
	auth.Use(testAuthMiddleware(userID))
	auth.GET("/me/threads", h.ListMine)
	auth.POST("/threads", h.Create)
	auth.PUT("/threads/:id", h.Update)
	auth.DELETE("/threads/:id", h.Delete)

	return r
}

func TestThreadList(t *testing.T) {
	repo := &fakeThreadRepo{
		listResult: []models.Thread{
			{Model: gorm.Model{ID: 1}, Title: "t1", UserID: 1},
		},
		countResult: 1,
	}
	r := newThreadRouter(repo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp dto.ThreadListResp
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestThreadListCursorInvalid(t *testing.T) {
	repo := &fakeThreadRepo{}
	r := newThreadRouter(repo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads?cursor=bad", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestThreadListCursorOK(t *testing.T) {
	ts := time.Unix(0, 123)
	repo := &fakeThreadRepo{
		listAfterResult: []models.Thread{
			{Model: gorm.Model{ID: 7, CreatedAt: ts}, Title: "t1", UserID: 1},
		},
	}
	r := newThreadRouter(repo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads?cursor=1_1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp dto.ThreadListResp
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

func TestThreadDetailNotFound(t *testing.T) {
	repo := &fakeThreadRepo{findResult: nil}
	r := newThreadRouter(repo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestThreadCreateUnauthorized(t *testing.T) {
	repo := &fakeThreadRepo{}
	r := newThreadRouter(repo, 0)

	body := `{"title":"t","content":"c"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestThreadCreateBadBody(t *testing.T) {
	repo := &fakeThreadRepo{}
	r := newThreadRouter(repo, 1)

	body := `{"title":"","content":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestThreadCreateRepoError(t *testing.T) {
	repo := &fakeThreadRepo{createErr: errors.New("boom")}
	r := newThreadRouter(repo, 1)

	body := `{"title":"t","content":"c"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

func TestThreadCreateOK(t *testing.T) {
	repo := &fakeThreadRepo{}
	r := newThreadRouter(repo, 1)

	body := `{"title":"t","content":"c"}`
	req := httptest.NewRequest(http.MethodPost, "/api/threads", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}

	var resp dto.ThreadDetailResp
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.UserID != 1 || resp.Title != "t" || resp.Content != "c" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestThreadListMineUnauthorized(t *testing.T) {
	repo := &fakeThreadRepo{}
	r := newThreadRouter(repo, 0)

	req := httptest.NewRequest(http.MethodGet, "/api/me/threads", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestThreadListMineOK(t *testing.T) {
	repo := &fakeThreadRepo{
		listResult: []models.Thread{
			{Model: gorm.Model{ID: 1}, Title: "t1", UserID: 1},
		},
		countResult: 1,
	}
	r := newThreadRouter(repo, 1)

	req := httptest.NewRequest(http.MethodGet, "/api/me/threads?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp dto.ThreadListResp
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestThreadDetailBadParam(t *testing.T) {
	repo := &fakeThreadRepo{}
	r := newThreadRouter(repo, 0)

	req := httptest.NewRequest(http.MethodGet, "/threads/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestThreadUpdateNotFound(t *testing.T) {
	repo := &fakeThreadRepo{findResult: nil}
	r := newThreadRouter(repo, 1)

	body := `{"title":"t","content":"c"}`
	req := httptest.NewRequest(http.MethodPut, "/api/threads/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestThreadUpdateForbidden(t *testing.T) {
	repo := &fakeThreadRepo{
		findResult: &models.Thread{Model: gorm.Model{ID: 1}, UserID: 2},
	}
	r := newThreadRouter(repo, 1)

	body := `{"title":"t","content":"c"}`
	req := httptest.NewRequest(http.MethodPut, "/api/threads/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusForbidden, w.Code, w.Body.String())
	}
}

func TestThreadDeleteNotFound(t *testing.T) {
	repo := &fakeThreadRepo{findResult: nil}
	r := newThreadRouter(repo, 1)

	req := httptest.NewRequest(http.MethodDelete, "/api/threads/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusNotFound, w.Code, w.Body.String())
	}
}

func TestThreadDeleteForbidden(t *testing.T) {
	repo := &fakeThreadRepo{
		findResult: &models.Thread{Model: gorm.Model{ID: 1}, UserID: 2},
	}
	r := newThreadRouter(repo, 1)

	req := httptest.NewRequest(http.MethodDelete, "/api/threads/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusForbidden, w.Code, w.Body.String())
	}
}
