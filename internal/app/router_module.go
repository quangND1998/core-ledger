package app

import (
	"context"
	config "core-ledger/configs"
	"core-ledger/internal/module/transactions"
	"core-ledger/model/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	// "core-ledger/internal/module/accounts"
	// "core-ledger/internal/module/customers"
	// "core-ledger/internal/module/payments"
	// ... import thêm module khác
)

// RouterParams contains all dependencies needed for route setup
type RouterParams struct {
	fx.In

	Router             *gin.Engine
	Lifecycle          fx.Lifecycle
	TransactionHandler transactions.TransactionHandler
	// Add more handlers here as needed:
	// UserHandler    *handler.UserHandler
	// OrderHandler   *handler.OrderHandler
	// AccountHandler *handler.AccountHandler
}

// NewRouter is the fx provider that constructs the Gin engine
func NewRouter() *gin.Engine {
	// use default with logger and recovery middleware
	router := gin.New()
	return router
}

func SetupAllRoutes(params RouterParams) {
	api := params.Router.Group("/api/v2")

	// Health route
	params.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, dto.PreResponse{
			Data: gin.H{
				"status":  "success",
				"message": "Wealify API is running",
			}},
		)
	})

	// Protected routes (with middleware)
	// Option 1: Apply middleware to entire protected group
	protected := api.Group("")
	// protected.Use(authMiddleware, loggingMiddleware) // Uncomment when you have middleware

	// Option 2: Apply middleware per module (more flexible)
	// transactions.SetupRoutes(protected, params.TransactionHandler, authMiddleware, loggingMiddleware)

	// Gọi SetupRoutes() từng module
	// Without middleware:
	transactions.SetupRoutes(protected, params.TransactionHandler)
	// With middleware (example):
	// transactions.SetupRoutes(protected, params.TransactionHandler, transactions.AuthMiddleware(), transactions.LoggingMiddleware())

	// accounts.SetupRoutes(protected, params.AccountHandler)
	// customers.SetupRoutes(protected, params.CustomerHandler)
	// payments.SetupRoutes(protected, params.PaymentHandler)
}

// StartHTTPServer starts Gin using Fx lifecycle (non-blocking) and graceful shutdown
func StartHTTPServer(params RouterParams) {
	// CORS middleware
	params.Router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	port := config.GetConfig().Common.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: params.Router,
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logrus.Printf("Server starting on port %s", port)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logrus.Errorf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logrus.Printf("Server shutting down")
			return srv.Shutdown(ctx)
		},
	})
}

var RouterModule = fx.Module("router",
	fx.Provide(NewRouter),
	fx.Invoke(SetupAllRoutes, StartHTTPServer),
)
