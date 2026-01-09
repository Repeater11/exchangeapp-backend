package middleware

import (
	"encoding/json"
	"exchangeapp/pkg/jwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newAuthRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/me", Auth(secret), func(c *gin.Context) {
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		c.JSON(http.StatusOK, gin.H{
			"id":       userID,
			"username": username,
		})
	})
	return r
}

func TestAuthMissingHeader(t *testing.T) {
	r := newAuthRouter("secret")
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestAuthInvalidToken(t *testing.T) {
	r := newAuthRouter("secret")
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer bad")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}

func TestAuthValidToken(t *testing.T) {
	r := newAuthRouter("secret")
	tokenStr, err := jwt.GenerateToken(1, "alice", "secret", 60)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if resp.ID != 1 || resp.Username != "alice" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
