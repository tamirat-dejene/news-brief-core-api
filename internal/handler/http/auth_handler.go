package http

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
	"os"
)

type AuthHandler struct {
	UserUseCase      contract.IUserUseCase
	BaseURL          string
	config           contract.IConfigProvider
	jwtService       contract.IJWTService
}

func NewAuthHandler(uc contract.IUserUseCase, baseURL string, config contract.IConfigProvider, jwtSvc contract.IJWTService) *AuthHandler {
	return &AuthHandler{
		UserUseCase:     uc,
		BaseURL:         baseURL,
		config:          config,
		jwtService:      jwtSvc,
	}
}

type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *AuthHandler) googleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  h.BaseURL + "/api/v1/auth/google/callback",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func (h *AuthHandler) cookieParams() (domain string, secure bool) {
	u, err := url.Parse(h.BaseURL)
	if err != nil {
		return "", false
	}
	secure = u.Scheme == "https"
	return u.Hostname(), secure
}

func (h *AuthHandler) HandleGoogleLogin(ctx *gin.Context) {
	// Generate OAuth state with platform information
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	oauthStateString := base64.URLEncoding.EncodeToString(b)
	ctx.SetCookie("oauthState", oauthStateString, 300, "/", "", false, true)

	url := h.googleOauthConfig().AuthCodeURL(oauthStateString)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) HandleGoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	cookieState, err := ctx.Cookie("oauthState")

	if err != nil || state != cookieState {
		ctx.String(http.StatusUnauthorized, "invalid CSRF state token\n")
		return
	}
	cookieSecure := os.Getenv("OAUTH2_SET_COOKIE_SECURE")
	ctx.SetCookie("oauthState", "", -1, "/", "", cookieSecure == "true", true)

	code := ctx.Query("code")
	if code == "" {
		ctx.String(http.StatusBadRequest, "authorization code not provided")
		return
	}

	requestCtx := ctx.Request.Context()

	token, err := h.googleOauthConfig().Exchange(requestCtx, code)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("failed to exchange authorization for token: %v\n", err))
		return
	}

	client := h.googleOauthConfig().Client(requestCtx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get user info: %v", err))
		return
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("Failed to decode user info: %v\n", err))
		return
	}
	fullname := userInfo.Name

	accessToken, refreshToken, err := h.UserUseCase.LoginWithOAuth(requestCtx, fullname, userInfo.Email)

	accessToken, refreshToken, err := h.UserUseCase.LoginWithOAuth(requestCtx, fullName, userInfo.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("failed to login with OAuth: %v\n", err))
		return
	}

	// Platform is already retrieved from cookie above
	
	// Handle mobile app flow
	if platform == "mobile" {
		// For mobile, return JSON response with tokens
		ctx.JSON(http.StatusOK, gin.H{
			"message":       "login successful",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"platform":      "mobile",
		})
		return
	}

	// Handle web flow - redirect to frontend
	frontendURL := h.config.GetFrontendBaseURL()
	
	if frontendURL != "" {
		// Parse frontend URL and add tokens as query parameters
		u, err := url.Parse(frontendURL)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Invalid frontend URL configuration")
			return
		}
		
		// Set the path to the auth success page
		u.Path = "/auth/success"
		
		// Add tokens as query parameters (more secure than fragments)
		query := u.Query()
		query.Set("access_token", accessToken)
		query.Set("refresh_token", refreshToken)
		
		// Add user ID if we can parse the token
		if claims, err := h.jwtService.ParseAccessToken(accessToken); err == nil {
			query.Set("user_id", claims.UserID)
		}
		
		u.RawQuery = query.Encode()
		
		// Redirect to frontend with tokens
		ctx.Redirect(http.StatusFound, u.String())
		return
	}

	// Fallback: return JSON if no frontend URL is configured
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "login successful",

		"access token":  accessToken,
		"refresh token": refreshToken,
	})
}
