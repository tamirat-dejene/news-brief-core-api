package http

import (
	"time"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	"github.com/RealEskalate/G6-NewsBrief/internal/handler/http/middleware"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	userHandler         *UserHandler
	emailHandler        *EmailHandler
	authHandler         *AuthHandler
	sourceHandler       *SourceHandler
	topicHandler        *TopicHandler
	subscriptionHandler *SubscriptionHandler
	userUsecase         contract.IUserUseCase
	jwtService          contract.IJWTService
}

func NewRouter(userUsecase contract.IUserUseCase, emailVerUC contract.IEmailVerificationUC, userRepo contract.IUserRepository, tokenRepo contract.ITokenRepository, hasher contract.IHasher, jwtService contract.IJWTService, mailService contract.IEmailService, logger contract.IAppLogger, config contract.IConfigProvider, validator contract.IValidator, uuidGen contract.IUUIDGenerator, randomGen contract.IRandomGenerator, sourceUC contract.ISourceUsecase, topicUC contract.ITopicUsecase, subscriptionUC contract.ISubscriptionUsecase) *Router {
	baseURL := config.GetAppBaseURL()
	return &Router{
		userHandler:         NewUserHandler(userUsecase),
		emailHandler:        NewEmailHandler(emailVerUC, userRepo),
		authHandler:         NewAuthHandler(userUsecase, baseURL),
		sourceHandler:       NewSourceHandler(sourceUC, uuidGen),
		topicHandler:        NewTopicHandler(topicUC, uuidGen),
		subscriptionHandler: NewSubscriptionHandler(subscriptionUC),
		jwtService:          jwtService,
		userUsecase:         userUsecase,
	}
}

func (r *Router) SetupRoutes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// rate limiter configuration
	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	lmt.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"})
	lmt.SetMessage("Too many requests, please try again later.")
	router.Use(middleware.RateLimiter(lmt))

	// API docs
	RegisterDocsRoutes(router)

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.userHandler.CreateUser)
		auth.POST("/login", r.userHandler.Login)
		auth.GET("/verify-email", r.emailHandler.HandleVerifyEmailToken)
		auth.POST("/forgot-password", r.userHandler.ForgotPassword)
		auth.POST("/reset-password", r.userHandler.ResetPassword)
		auth.POST("/refresh-token", r.userHandler.RefreshToken)
		auth.POST("/request-verification-email", r.emailHandler.HandleRequestEmailVerification)

		// Google OAuth endpoints
		auth.GET("/google/login", r.authHandler.HandleGoogleLogin)
		auth.GET("/google/callback", r.authHandler.HandleGoogleCallback)
	}

	// Public user routes
	users := v1.Group("/users")
	{
		users.GET("/profile/:id", r.userHandler.GetUser)
	}

	// Protected routes (authentication required)
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleWare(r.jwtService, r.userUsecase))
	{
		// user routes
		protected.GET("/me", r.userHandler.GetCurrentUser)
		protected.PUT("/me", r.userHandler.UpdateUser)
		protected.GET("/me/subscriptions", r.subscriptionHandler.GetSubscriptions)
		protected.POST("/me/subscriptions", r.subscriptionHandler.AddSubscription)
		protected.DELETE("/me/subscriptions/:source_slug", r.subscriptionHandler.RemoveSubscription)
		// admin routes
		protected.POST("/topics", r.topicHandler.CreateTopic)
		// protected.POST("/topics/:id", r.topicHandler.UpdateTopic)
		// protected.DELETE("/topics/:id", r.topicHandler.DeleteTopic)
		protected.POST("/sources", r.sourceHandler.CreateSource)
		// protected.PUT("/sources/:id", r.sourceHandler.UpdateSource)
		// protected.DELETE("/sources/:id", r.sourceHandler.DeleteSource)
	}

	// Logout route (no authentication required just accept the refresh token from the request body and invalidate the user session)
	v1.POST("/logout", r.userHandler.Logout)
}
