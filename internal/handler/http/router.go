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
		topicHandler:        NewTopicHandler(topicUC, userUsecase, uuidGen),
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

	// router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// router.GET("/api/v1/metrics", gin.WrapH(promhttp.Handler()))

	// API v1 routes
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", r.userHandler.CreateUser)
		auth.POST("/login", r.userHandler.Login)
		auth.POST("/refresh-token", r.userHandler.RefreshToken)
		auth.GET("/verify-email", r.emailHandler.HandleVerifyEmailToken)
		auth.POST("/forgot-password", r.userHandler.ForgotPassword)
		auth.POST("/reset-password", r.userHandler.ResetPassword)
		auth.POST("/request-verification-email", r.emailHandler.HandleRequestEmailVerification)
		// Google OAuth endpoints
		auth.GET("/google/login", r.authHandler.HandleGoogleLogin)
		auth.GET("/google/callback", r.authHandler.HandleGoogleCallback)
	}

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleWare(r.jwtService, r.userUsecase))
	{
		// admin routes
		admin.POST("/create-topics", r.topicHandler.CreateTopic)
		// admin.POST("/topics/:id", r.topicHandler.UpdateTopic)
		// admin.DELETE("/topics/:id", r.topicHandler.DeleteTopic)
		admin.POST("/create-sources", r.sourceHandler.CreateSource)
		// admin.PUT("/sources/:id", r.sourceHandler.UpdateSource)
		// admin.DELETE("/sources/:id", r.sourceHandler.DeleteSource)
	}

	// user profile routes (authentication required)
	userProfile := v1.Group("/me")
	userProfile.Use(middleware.AuthMiddleWare(r.jwtService, r.userUsecase))
	{
		// user routes
		userProfile.GET("", r.userHandler.GetCurrentUser)
		userProfile.PUT("", r.userHandler.UpdateUser)
		userProfile.GET("/subscriptions", r.subscriptionHandler.GetSubscriptions)
		userProfile.POST("/subscriptions", r.subscriptionHandler.AddSubscription)
		userProfile.DELETE("/subscriptions/:source_slug", r.subscriptionHandler.RemoveSubscription)
		userProfile.POST("/topics", r.topicHandler.SubscribeTopic)
		userProfile.GET("/subscribed-topics", r.topicHandler.GetUserSubscribedTopics)
	}
	// public api
	public := v1.Group("")
	{
		public.GET("/topics", r.topicHandler.GetTopics)
		public.GET("/sources", r.sourceHandler.GetSources)
	}
	// Logout route (no authentication required just accept the refresh token from the request body and invalidate the user session)
	v1.POST("/logout", r.userHandler.Logout)
}
