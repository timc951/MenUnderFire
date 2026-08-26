package main

import (
	"net/http"
	"os"
	"time"

	"menunderfire/internal/config"
	"menunderfire/internal/database"
	"menunderfire/internal/handlers"
	"menunderfire/internal/logger"
	"menunderfire/internal/middleware"
	"menunderfire/internal/repositories/postgres"
	"menunderfire/internal/routes"
	"menunderfire/internal/services"

	"github.com/joho/godotenv"
)

// Version is set at build time via -ldflags
var Version = "dev"

func main() {
	// Initialize logger based on environment
	env := os.Getenv("APP_ENV")
	if env == "production" {
		logger.InitProduction()
	} else {
		logger.InitDefault()
	}

	log := logger.Log

	log.Info().Msg("========================================")
	log.Info().Str("version", Version).Msg("   MenUnderFireApp Go Backend Starting")
	log.Info().Msg("========================================")

	// Load .env file
	log.Debug().Msg("Loading .env file...")
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg(".env file not found, using environment variables")
	} else {
		log.Debug().Msg(".env file loaded successfully")
	}

	// Load configuration
	log.Debug().Msg("Loading configuration...")
	cfg := config.Load()
	log.Info().
		Str("port", cfg.ServerPort).
		Str("db_host", cfg.DBHost).
		Str("db_port", cfg.DBPort).
		Str("db_name", cfg.DBName).
		Str("auth_domain", cfg.AuthDomain).
		Str("auth_issuer", cfg.AuthIssuer).
		Msg("Configuration loaded")

	// Connect to database
	log.Info().Msg("Connecting to database...")
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connection established")
	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database")
		}
	}()

	// Run database migrations
	log.Info().Msg("Running database migrations...")
	migrationsApplied, err := database.RunMigrations(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}
	if migrationsApplied > 0 {
		log.Info().Int("count", migrationsApplied).Msg("Database migrations applied")
	} else {
		log.Info().Msg("Database schema is up to date")
	}

	// ===== Initialize Repositories =====
	log.Debug().Msg("Initializing repositories...")
	userRepo := postgres.NewUserRepository(db)
	orgRepo := postgres.NewOrganizationRepository(db)
	groupRepo := postgres.NewGroupRepository(db)
	apiKeyRepo := postgres.NewAPIKeyRepository(db)
	formRepo := postgres.NewFormRepository(db)
	formFieldRepo := postgres.NewFormFieldRepository(db)
	formAnswerRepo := postgres.NewFormAnswerRepository(db)
	groupMessageRepo := postgres.NewGroupMessageRepository(db)
	invitationRepo := postgres.NewInvitationRepository(db)
	reportRepo := postgres.NewReportRepository(db)
	sitePageRepo := postgres.NewSitePageRepository(db)
	pageDraftRepo := postgres.NewPageDraftRepository(db)
	feedbackRepo := postgres.NewFeedbackRepository(db)
	pageHitRepo := postgres.NewPageHitRepository(db)
	log.Debug().Msg("Repositories initialized")

	// ===== Initialize Services =====
	log.Debug().Msg("Initializing services...")
	userService := services.NewUserService(userRepo, cfg.AgreementHMACSecret)
	orgService := services.NewOrganizationService(orgRepo, groupRepo, userRepo)
	groupService := services.NewGroupService(groupRepo, orgRepo, userRepo)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	formService := services.NewFormService(formRepo, formFieldRepo, formAnswerRepo, orgRepo, userRepo)
	groupMessageService := services.NewGroupMessageService(groupMessageRepo, groupRepo, groupService, orgRepo, userRepo, formRepo, formAnswerRepo)
	invitationService := services.NewInvitationService(invitationRepo, userRepo, orgRepo, groupRepo, groupService)
	reportService := services.NewReportService(reportRepo, groupRepo)
	sitePageService := services.NewSitePageService(sitePageRepo, pageDraftRepo, userRepo)
	dashboardService := services.NewDashboardService(userService, orgRepo, groupRepo)
	feedbackService := services.NewFeedbackService(feedbackRepo, userRepo)
	pageHitService := services.NewPageHitService(pageHitRepo, userRepo)
	log.Debug().Msg("Services initialized")

	// ===== Initialize Middleware =====
	log.Debug().Msg("Initializing auth middleware...")
	authMiddleware := middleware.NewAuthMiddleware(cfg, userRepo)
	log.Debug().Msg("Auth middleware initialized")

	// ===== Initialize Handlers =====
	log.Debug().Msg("Initializing handlers...")
	userHandler := handlers.NewUserHandler(userService)
	groupHandler := handlers.NewGroupHandler(groupService)
	groupMessageHandler := handlers.NewGroupMessageHandler(groupMessageService)
	reportHandler := handlers.NewReportHandler(reportService, groupService)
	organizationHandler := handlers.NewOrganizationHandler(orgService, userService)
	invitationHandler := handlers.NewInvitationHandler(invitationService)
	formHandler := handlers.NewFormHandler(formService)
	sitePageHandler := handlers.NewSitePageHandler(sitePageService, userService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService, userService)
	apiKeyHandler := handlers.NewAPIKeyHandler(apiKeyService)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackService, userService)
	pageHitHandler := handlers.NewPageHitHandler(pageHitService, cfg.HitToken)
	log.Debug().Msg("Handlers initialized")

	// ===== Setup Routes =====
	log.Debug().Msg("Setting up routes...")
	router := routes.Setup(
		userHandler,
		groupHandler,
		groupMessageHandler,
		reportHandler,
		organizationHandler,
		invitationHandler,
		formHandler,
		sitePageHandler,
		dashboardHandler,
		apiKeyHandler,
		feedbackHandler,
		pageHitHandler,
		authMiddleware,
		cfg.CORSOrigin,
		cfg.TrustProxyHeaders,
	)
	log.Debug().Msg("Routes configured")

	// ===== Start Server =====
	addr := ":" + cfg.ServerPort
	log.Info().Msg("========================================")
	log.Info().Str("address", "http://localhost"+addr).Msg("Server ready")
	log.Info().Msg("========================================")

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
