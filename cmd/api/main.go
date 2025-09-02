package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/RealEskalate/G6-NewsBrief/internal/domain/contract"
	handlerHttp "github.com/RealEskalate/G6-NewsBrief/internal/handler/http"
	"github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/config"
	database "github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/database"
	"github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/external_services"
	"github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/jwt"
	"github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/logger"
	passwordservice "github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/password_service"
	randomgenerator "github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/random_generator"
	"github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/repository/mongodb"
	"github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/seeder"
	"github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/uuidgen"
	"github.com/RealEskalate/G6-NewsBrief/internal/infrastructure/validator"
	"github.com/RealEskalate/G6-NewsBrief/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func adminSeeder(userUsecase contract.IUserUseCase, userRepo contract.IUserRepository) {
	// Add a -seed flag to run seeder before starting the server
	seed := flag.Bool("seed", true, "run database seeder and exit")
	flag.Parse()

	// Optionally allow seeding via env (useful on Render/CI)
	seedOnStart := os.Getenv("SEED_ON_START") == "true"

	if *seed || seedOnStart {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		adminEmail := os.Getenv("SEED_ADMIN_EMAIL")
		if adminEmail == "" {
			adminEmail = "admin@newsbrief.local"
		}
		adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")
		if adminPassword == "" {
			adminPassword = "ChangeMe123!"
		}

		if err := seeder.SeedAdminUsingUC(ctx, userUsecase, userRepo, adminEmail, adminPassword); err != nil {
			log.Fatalf("seeding failed: %v", err)
		}
		log.Println("seeding completed")
		// Exit if seeding-only mode
		if *seed {
			return
		}
	}
}

func adminSeeder(userUsecase contract.IUserUseCase, userRepo contract.IUserRepository) {
	// Add a -seed flag to run seeder before starting the server
	seed := flag.Bool("seed", true, "run database seeder and exit")
	flag.Parse()

	// Optionally allow seeding via env (useful on Render/CI)
	seedOnStart := os.Getenv("SEED_ON_START") == "true"

	if *seed || seedOnStart {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		adminEmail := os.Getenv("SEED_ADMIN_EMAIL")
		if adminEmail == "" {
			adminEmail = "admin@newsbrief.local"
		}
		adminPassword := os.Getenv("SEED_ADMIN_PASSWORD")
		if adminPassword == "" {
			adminPassword = "ChangeMe123!"
		}

		if err := seeder.SeedAdminUsingUC(ctx, userUsecase, userRepo, adminEmail, adminPassword); err != nil {
			log.Fatalf("seeding failed: %v", err)
		}
		log.Println("seeding completed")
		// Exit if seeding-only mode
		if *seed {
			return
		}
	}
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get MongoDB URI and DB name from environment
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGODB_DB_NAME environment variable not set")
	}

	// Establish MongoDB connection
	mongoClient, err := database.NewMongoDBClient(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect()

	// Initialize email service (SendGrid)
	sendGridAPIKey := os.Getenv("SENDGRID_API_KEY")
	sendFrom := os.Getenv("EMAIL_FROM")
	sendFromName := os.Getenv("EMAIL_FROM_NAME")

	// Register custom validators
	validator.RegisterCustomValidators()

	// Dependency Injection: Repositories
	userCollection := mongoClient.Client.Database(dbName).Collection("users")
	userRepo := mongodb.NewUserRepository(userCollection)
	tokenRepo := mongodb.NewTokenRepository(mongoClient.Client.Database(dbName).Collection("tokens"))
	topicRepo := mongodb.NewTopicRepository(mongoClient.Client.Database(dbName).Collection("topics"))
	sourceRepo := mongodb.NewSourceRepository(mongoClient.Client.Database(dbName).Collection("sources"))
	// Dependency Injection: Services
	hasher := passwordservice.NewHasher()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	jwtManager := jwt.NewJWTManager(jwtSecret)
	jwtService := jwt.NewJWTService(jwtManager)
	appLogger := logger.NewStdLogger()
	mailService := external_services.NewEmailService(sendGridAPIKey, sendFrom, sendFromName)
	randomGenerator := randomgenerator.NewRandomGenerator()
	appValidator := validator.NewValidator()
	uuidGenerator := uuidgen.NewGenerator()
	appConfig := config.NewConfig()
	// Dependency Injection: Usecases
	emailUsecase := usecase.NewEmailVerificationUseCase(tokenRepo, userRepo, mailService, randomGenerator, uuidGenerator, baseURL)
	userUsecase := usecase.NewUserUsecase(userRepo, tokenRepo, topicRepo, emailUsecase, hasher, jwtService, mailService, appLogger, appConfig, appValidator, uuidGenerator, randomGenerator)
	topicUsecase := usecase.NewTopicUsecase(topicRepo)
	sourceUsecase := usecase.NewSourceUsecase(sourceRepo)
	subscriptionUsecase := usecase.NewSubscriptionUsecase(userRepo, sourceRepo)
	// Pass Prometheus metrics to handlers or usecases as needed (import from metrics package)

	//---------------------- Admin seeder-------------------------------------
	adminSeeder(userUsecase, userRepo)
	//---------------------- end of admin seeder-------------------------------------

	// Setup API routes
	appRouter := handlerHttp.NewRouter(
		userUsecase, emailUsecase,
		userRepo, tokenRepo, hasher, jwtService, mailService,
		appLogger, appConfig, appValidator, uuidGenerator, randomGenerator, sourceUsecase, topicUsecase, subscriptionUsecase,
	)

	// Initialize Gin router
	router := gin.Default()

	appRouter.SetupRoutes(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
