package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"gin-demo/handler"
	"gin-demo/middleware"
	"gin-demo/mocks"
	"gin-demo/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler *handler.Handler) *gin.Engine {
	r := gin.Default()
	r.POST("/api/login", handler.Login)
	r.POST("/api/logout", handler.Logout)
	r.POST("/api/users", handler.Register)
	r.GET("/api/users", handler.GetUsers)
	r.GET("/api/users/me", handler.GetAuthenticatedUser)
	r.GET("/api/user/:id", handler.GetUserById)
	return r
}

func TestLogin_Success(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	reqBody := model.UserLoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	resp := &model.UserLoginResponse{JWTToken: "jwt-token"}

	mockService.On("Login", &reqBody).Return(resp, nil)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := setupRouter(handler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Logged in successfully")
	mockService.AssertExpectations(t)
}

func TestLogin_InvalidRequest(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString("{invalid-json}"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := setupRouter(handler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_ServiceError(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	reqBody := model.UserLoginRequest{
		Email:    "fail@example.com",
		Password: "wrong",
	}
	mockService.On("Login", &reqBody).Return(nil, errors.New("invalid credentials"))

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := setupRouter(handler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
}

func TestLogout_Success(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("claims", &middleware.Claims{Email: "test@example.com"})

	mockService.On("Logout", &model.UserLogoutRequest{Email: "test@example.com"}).
		Return(&model.UserLogoutResponse{Message: "Logout success"}, nil)

	handler.Logout(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Logout success")
}

func TestLogout_Error(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("claims", &middleware.Claims{Email: "test@example.com"})
	mockService.On("Logout", &model.UserLogoutRequest{Email: "test@example.com"}).
		Return(nil, errors.New("logout failed"))

	handler.Logout(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRegister_Success(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	user := model.User{
		Email:    "new@example.com",
		Password: "password",
	}
	mockService.On("Register", &user).Return(nil)

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := setupRouter(handler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Registered successfully")
}

func TestRegister_Error(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	user := model.User{Email: "bad@example.com"}
	mockService.On("Register", &user).Return(errors.New("email exists"))

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := setupRouter(handler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestGetUsers_Success(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	users := []model.UserData{{Email: "a@example.com"}, {Email: "b@example.com"}}
	mockService.On("GetUsers", "test@example.com").Return(users)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("claims", &middleware.Claims{Email: "test@example.com"})

	handler.GetUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "a@example.com")
}

func TestGetAuthenticatedUser_Success(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	user := &model.UserData{Email: "me@example.com"}
	mockService.On("GetUserByEmail", "me@example.com").Return(user, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("claims", &middleware.Claims{Email: "me@example.com"})

	handler.GetAuthenticatedUser(c)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "me@example.com")
}

func TestGetUserById_Success(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	user := &model.UserData{
		Username:  "tester",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	}
	mockService.On("GetUserById", "1").Return(user, nil)

	req, _ := http.NewRequest(http.MethodGet, "/api/user/1", nil)
	w := httptest.NewRecorder()
	r := setupRouter(handler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "tester")
}

func TestGetUserById_InvalidID(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/api/user/abc", nil)
	w := httptest.NewRecorder()
	r := setupRouter(handler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserById_NotFound(t *testing.T) {
	mockService := new(mocks.UserService)
	handler := handler.NewHandler(mockService)

	mockService.On("GetUserById", "99").Return(nil, errors.New("user not found"))

	req, _ := http.NewRequest(http.MethodGet, "/api/user/99", nil)
	w := httptest.NewRecorder()
	r := setupRouter(handler)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
