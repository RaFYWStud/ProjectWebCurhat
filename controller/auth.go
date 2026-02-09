package controller

import (
	"net/http"

	"projectwebcurhat/config/middleware"
	"projectwebcurhat/contract"
	"projectwebcurhat/dto"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *contract.Service
}

func (a *AuthController) GetPrefix() string {
	return "/auth"
}

func (a *AuthController) InitService(service *contract.Service) {
	a.service = service
}

func (a *AuthController) InitRoute(app *gin.RouterGroup) {
	app.POST("/register", a.Register)
	app.POST("/login", a.Login)
	app.GET("/profile", middleware.AuthMiddleware(), a.GetProfile)
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.RegisterRequest true "Register payload"
// @Success 201 {object} dto.AuthResponse
// @Router /auth/register [post]
func (a *AuthController) Register(ctx *gin.Context) {
	var payload dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := a.service.Auth.Register(&payload)
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful",
		"data":    result,
	})
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login payload"
// @Success 200 {object} dto.AuthResponse
// @Router /auth/login [post]
func (a *AuthController) Login(ctx *gin.Context) {
	var payload dto.LoginRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := a.service.Auth.Login(&payload)
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    result,
	})
}

// GetProfile godoc
// @Summary Get current user profile
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserProfile
// @Router /auth/profile [get]
func (a *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	profile, err := a.service.Auth.GetProfile(userID.(int))
	if err != nil {
		HandlerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}
