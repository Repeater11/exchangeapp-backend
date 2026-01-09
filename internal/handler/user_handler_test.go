package handler

import (
	"encoding/json"
	"errors"
	"exchangeapp/internal/config"
	"exchangeapp/internal/dto"
	"exchangeapp/internal/models"
	"exchangeapp/internal/repository"
	"exchangeapp/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type fakeUserRepo struct {
	createErr            error
	findByUsernameResult *models.User
	findByUsernameErr    error
	nextID               uint
}

func (f *fakeUserRepo) Create(u *models.User) error {
	if f.createErr != nil {
		return f.createErr
	}
	if f.nextID == 0 {
		f.nextID = 1
	}
	u.ID = f.nextID
	return nil
}

func (f *fakeUserRepo) FindByUsername(username string) (*models.User, error) {
	return f.findByUsernameResult, f.findByUsernameErr
}

func newUserRouter(repo repository.UserRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := service.NewUserService(repo, config.JWTConfig{
		Secret:        "test",
		ExpireMinutes: 60,
	})
	h := NewUserHandler(svc)

	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	return r
}

func TestRegister(t *testing.T) {
	cases := []struct {
		name     string
		repo     repository.UserRepository
		body     string
		wantCode int
	}{
		{
			name:     "ok",
			repo:     &fakeUserRepo{},
			body:     `{"username":"alice","password":"pass123"}`,
			wantCode: http.StatusCreated,
		},
		{
			name:     "bad_body",
			repo:     &fakeUserRepo{},
			body:     `{"username":"a","password":"123"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "conflict",
			repo:     &fakeUserRepo{createErr: repository.ErrUserExists},
			body:     `{"username":"alice","password":"pass123"}`,
			wantCode: http.StatusConflict,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newUserRouter(c.repo)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(c.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != c.wantCode {
				t.Fatalf("expected %d, got %d, body=%s", c.wantCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestLogin(t *testing.T) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password failed: %v", err)
	}

	cases := []struct {
		name      string
		repo      repository.UserRepository
		body      string
		wantCode  int
		wantToken bool
	}{
		{
			name: "ok",
			repo: &fakeUserRepo{
				findByUsernameResult: &models.User{
					Model:    gorm.Model{ID: 1},
					Username: "alice",
					Password: string(hashed),
				},
			},
			body:      `{"username":"alice","password":"pass123"}`,
			wantCode:  http.StatusOK,
			wantToken: true,
		},
		{
			name:      "bad_body",
			repo:      &fakeUserRepo{},
			body:      `{"username":"alice","password":"123"}`,
			wantCode:  http.StatusBadRequest,
			wantToken: false,
		},
		{
			name:      "user_not_found",
			repo:      &fakeUserRepo{findByUsernameResult: nil},
			body:      `{"username":"alice","password":"pass123"}`,
			wantCode:  http.StatusUnauthorized,
			wantToken: false,
		},
		{
			name:      "repo_error",
			repo:      &fakeUserRepo{findByUsernameErr: errors.New("boom")},
			body:      `{"username":"alice","password":"pass123"}`,
			wantCode:  http.StatusInternalServerError,
			wantToken: false,
		},
		{
			name: "wrong_password",
			repo: &fakeUserRepo{
				findByUsernameResult: &models.User{
					Model:    gorm.Model{ID: 1},
					Username: "alice",
					Password: string(hashed),
				},
			},
			body:      `{"username":"alice","password":"pass124"}`,
			wantCode:  http.StatusUnauthorized,
			wantToken: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := newUserRouter(c.repo)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(c.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != c.wantCode {
				t.Fatalf("expected %d, got %d, body=%s", c.wantCode, w.Code, w.Body.String())
			}

			if c.wantToken {
				var resp dto.LoginResp
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("unmarshal failed: %v", err)
				}
				if resp.Token == "" {
					t.Fatalf("expected token, got empty")
				}
			}
		})
	}
}
