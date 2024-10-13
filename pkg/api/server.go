package api

import (
	"stream-app/pkg/api/handlers"
	"stream-app/pkg/api/middleware"
	"stream-app/pkg/config"
	"stream-app/pkg/logger"
	"stream-app/pkg/repositories"
	"stream-app/pkg/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	db     *gorm.DB
	logger *logger.Logger
}

func NewServer(cfg *config.Config, db *gorm.DB, logger *logger.Logger) *Server {

	router := gin.Default()
	router.LoadHTMLGlob("frontend/templates/*")
	server := &Server{
		router: router,
		cfg:    cfg,
		db:     db,
		logger: logger,
	}

	server.setupRoutes()
	return server

}

func (s *Server) setupRoutes() {

	// Middleware
	s.router.Use(gin.Recovery())
	s.router.Use(middleware.LoggerMiddleware(s.logger))
	s.router.Use(middleware.CORSMiddleware())

	streamRepo := repositories.NewStreamRepository(s.db)
	streamService := services.NewStreamService(streamRepo)

	// Routes
	authHandler := handlers.NewAuthHandler(s.db, s.cfg, s.logger)
	streamHandler := handlers.NewStreamHandler(streamService)

	// Routes whithout authentication
	s.router.GET("/", handlers.IndexHandler)
	s.router.GET("/login", handlers.LoginPageHandler)
	s.router.POST("/login", authHandler.Login)
	s.router.POST("/register", authHandler.Register)

	// Routes with authentication
	authorized := s.router.Group("/")
	authorized.Use(middleware.AuthMiddleware(s.cfg.JWT.SecretKey))
	{
		authorized.GET("/streams", streamHandler.GetAllStreamsHandler)
		authorized.POST("/streams", streamHandler.CreateStreamHandler)
		authorized.GET("/streams/:id", streamHandler.GetStreamHandler)
		authorized.DELETE("/streams/:id", streamHandler.DeleteStreamHandler)

	}

}

func (s *Server) Run() error {
	address := s.cfg.Server.Address
	s.logger.SugaredLogger.Infof("Server is running on %s", address)
	return s.router.Run(address)
}
