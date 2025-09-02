package http

import (
	"net/http"
	"net/url"
	"time"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailVerificationUC contract.IEmailVerificationUC
	userRepository      contract.IUserRepository
	jwtService          contract.IJWTService
	tokenRepo           contract.ITokenRepository
	hasher              contract.IHasher
	config              contract.IConfigProvider
	uuidGen             contract.IUUIDGenerator
}

func NewEmailHandler(eu contract.IEmailVerificationUC, uc contract.IUserRepository, jwtSvc contract.IJWTService, tokenRepo contract.ITokenRepository, hasher contract.IHasher, cfg contract.IConfigProvider, uuidGen contract.IUUIDGenerator) *EmailHandler {
	return &EmailHandler{
		emailVerificationUC: eu,
		userRepository:      uc,
		jwtService:          jwtSvc,
		tokenRepo:           tokenRepo,
		hasher:              hasher,
		config:              cfg,
		uuidGen:             uuidGen,
	}
}

type requestEmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *EmailHandler) HandleRequestEmailVerification(ctx *gin.Context) {
	var req requestEmailVerificationRequest
	requestCtx := ctx.Request.Context()
	// parse the json req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// fetch user based on the req.email
	user, err := h.userRepository.GetUserByEmail(requestCtx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// send a email validation request
	if err = h.emailVerificationUC.RequestVerificationEmail(requestCtx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}
	// send a successfull message
	ctx.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

func (h *EmailHandler) HandleVerifyEmailToken(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()
	verifier := ctx.Query("verifier")
	plainToken := ctx.Query("token")

	if verifier == "" || plainToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing verifier or token"})
		return
	}

	// call the verify email token usecase
	user, err := h.emailVerificationUC.VerifyEmailToken(requestCtx, verifier, plainToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid token or expired token"})
		return
	}
	user.IsVerified = true
	// update the user
	if _, err := h.userRepository.UpdateUser(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Issue tokens immediately after successful verification
	accessToken, err := h.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}
	expiry := h.config.GetRefreshTokenExpiry()
	if expiry <= 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid refresh token expiry configuration"})
		return
	}
	tokenEntity := &entity.Token{
		ID:        h.uuidGen.NewUUID(),
		UserID:    user.ID,
		TokenType: entity.TokenTypeRefresh,
		TokenHash: h.hasher.HashString(refreshToken),
		ExpiresAt: time.Now().Add(expiry),
		CreatedAt: time.Now(),
		Revoke:    false,
	}
	if err := h.tokenRepo.CreateToken(requestCtx, tokenEntity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	// success response with tokens
	frontend := h.config.GetFrontendBaseURL()
	mobile := h.config.GetFrontendMobileBaseURL()
	platform := ctx.Query("platform") // optional hint: "mobile" or "web"

	// Force JSON for mobile platform
	if platform == "mobile" {
		ctx.JSON(http.StatusOK, gin.H{
			"message":       "Email verified successfully",
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
		return
	}

	redirectBase := frontend
	if platform == "mobile" && mobile != "" {
		redirectBase = mobile
	}
	if redirectBase != "" {
		u, _ := url.Parse(redirectBase)
		u.Path = "/auth/verified"
		fragment := url.Values{}
		fragment.Set("access_token", accessToken)
		fragment.Set("refresh_token", refreshToken)
		fragment.Set("user_id", user.ID)
		u.Fragment = fragment.Encode()
		ctx.Redirect(http.StatusFound, u.String())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Email verified successfully",
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
