package handler

import (
	err "gin-demo/errors"
	"gin-demo/middleware"
	"gin-demo/model"
	"gin-demo/redis_utils"
	services "gin-demo/service"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services.UserService
}

func NewHandler(userService services.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}
func (h *Handler) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if er := c.ShouldBindJSON(&req); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.InvalidReqData})
		return
	}

	res, er := h.UserService.Login(&req)
	if er != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": er.Error()})
		return
	}

	authHeader := "Bearer " + res.JWTToken

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    authHeader,
		Path:     "/",
		Domain:   redis_utils.AppConfig.Cookie.Domain,
		Expires:  time.Now().Add(8 * time.Hour),
		MaxAge:   8 * 3600,
		HttpOnly: false,
		Secure:   false,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"email":   req.Email,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	claims, exists := c.Get("claims")
	email := claims.(*middleware.Claims).Email
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.UserNotAuthenticated})
		return
	}

	logoutReq := &model.UserLogoutRequest{
		Email: email,
	}

	logoutRes, er := h.UserService.Logout(logoutReq)
	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": logoutRes.Message,
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req model.User
	if er := c.ShouldBindJSON(&req); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}
	err := h.UserService.Register(&req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

func (h *Handler) GetUsers(c *gin.Context) {
	claims, _ := c.Get("claims")
	email := claims.(*middleware.Claims).Email
	users := h.UserService.GetUsers(email)
	c.JSON(http.StatusOK, gin.H{
		"All registered users": users,
	})
}

func (h *Handler) GetAuthenticatedUser(c *gin.Context) {
	claims, _ := c.Get("claims")
	email := claims.(*middleware.Claims).Email
	user, err := h.UserService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"User data": user,
	})
}

func (h *Handler) GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	if matched := regexp.MustCompile(`^\d+$`).MatchString(idParam); !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format — only digits are allowed"})
		return
	}
	user, err := h.UserService.GetUserById(idParam)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userData := model.UserData{
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userData,
	})
}
