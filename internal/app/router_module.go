package app

import (
	"context"
	config "core-ledger/configs"
	coaaccount "core-ledger/internal/module/coaAccount"
	"core-ledger/internal/module/entries"
	"core-ledger/internal/module/excel"
	"core-ledger/internal/module/middleware"
	"core-ledger/internal/module/option"
	"core-ledger/internal/module/permission"
	"core-ledger/internal/module/role"
	"core-ledger/internal/module/ruleCategory"
	"core-ledger/internal/module/ruleValue"
	"core-ledger/internal/module/transactions"
	"core-ledger/internal/module/user"
	"core-ledger/model/dto"
	"core-ledger/pkg/repo"
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

	Router              *gin.Engine
	Lifecycle           fx.Lifecycle
	TransactionHandler  *transactions.TransactionHandler
	ExcelHandler        *excel.ExcelHandler
	CoaAccountHandler   *coaaccount.CoaAccountHandler
	EntriesHandler      *entries.EntriesHandler
	RuleCategoryHandler *ruleCategory.RuleCategoryHandler
	RuleValueHander     *ruleValue.RuleValueHandler
	OptionHandler       *option.OptionHandler
	PermissionHandler   *permission.PermissionHandler
	RoleHandler         *role.RoleHandler
	UserHandler         *user.UserHandler
	AuthHandler         *user.AuthHandler
	CoaRequestHandler   *coaaccount.RequestCoaAccountHandler
	UserRepo            repo.UserRepo

	// Add more handlers here as needed:
	// UserHandler    *handler.UserHandler
	// OrderHandler   *handler.OrderHandler
	// AccountHandler *handler.AccountHandler
}

// NewRouter is the fx provider that constructs the Gin engine
func NewRouter() *gin.Engine {
	// use default with logger and recovery middleware
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.LogRequest)
	router.Use(func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: false,
	// 	MaxAge:           12 * time.Hour,
	// }))
	// router.Use(middleware.RateLimitMiddleware())
	return router
}

func SetupAllRoutes(params RouterParams) {
	api := params.Router.Group("/api/v1")
	params.Router.GET("/", func(c *gin.Context) {
		c.JSON(200, dto.PreResponse{
			Data: gin.H{
				"status":  "success",
				"message": "Core Ledger API is running",
			}},
		)
	})
	// Health route
	params.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, dto.PreResponse{
			Data: gin.H{
				"status":  "success",
				"message": "Core ledger API is running",
			}},
		)
	})
	params.Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Route not found",
			"path":    c.Request.URL.Path,
			"method":  c.Request.Method,
			"message": "The requested endpoint does not exist.",
		})
	})
	// Protected routes (with middleware)
	// Option 1: Apply middleware to entire protected group
	protected := api.Group("")
	// protected.Use(middleware.RateLimitMiddleware())
	// protected.Use(authMiddleware, loggingMiddleware) // Uncomment when you have middleware

	// Option 2: Apply middleware per module (more flexible)
	// transactions.SetupRoutes(protected, params.TransactionHandler, authMiddleware, loggingMiddleware)

	// Gọi SetupRoutes() từng module
	// Without middleware:
	userAuthMiddleware := user.UserAuthMiddleware(params.UserRepo)
	transactions.SetupRoutes(protected, params.TransactionHandler)
	excel.SetupRoutes(protected, params.ExcelHandler)
	// coaaccount.SetupRoutes(protected, params.CoaAccountHandler, params.CoaRequestHandler, userAuthMiddleware)
	coaaccount.SetupRoutes(protected, params.CoaAccountHandler, params.CoaRequestHandler)
	entries.SetupRoutes(protected, params.EntriesHandler)
	ruleCategory.SetupRoutes(protected, params.RuleCategoryHandler)
	ruleValue.SetupRoutes(protected, params.RuleValueHander)
	option.SetupRoutes(protected, params.OptionHandler)
	permission.SetupRoutes(protected, params.PermissionHandler, userAuthMiddleware)
	role.SetupRoutes(protected, params.RoleHandler, userAuthMiddleware)

	// Auth routes (public - no authentication required)
	user.SetupAuthRoutes(protected, params.AuthHandler)

	// User routes with authentication middleware

	user.SetupRoutes(protected, params.UserHandler, userAuthMiddleware)

	// With middleware (example):
	// transactions.SetupRoutes(protected, params.TransactionHandler, transactions.AuthMiddleware(), transactions.LoggingMiddleware())
	// permission.SetupRoutes(protected, params.PermissionHandler, middleware.AuthMiddleware(), middleware.PermissionMiddleware("manage-permissions", "web"))

	// accounts.SetupRoutes(protected, params.AccountHandler)
	// customers.SetupRoutes(protected, params.CustomerHandler)
	// payments.SetupRoutes(protected, params.PaymentHandler)
}

// StartHTTPServer starts Gin using Fx lifecycle (non-blocking) and graceful shutdown
func StartHTTPServer(params RouterParams) {
	// CORS middleware
	// params.Router.Use(func(c *gin.Context) {

	// 	c.Header("Access-Control-Allow-Origin", "*")
	// 	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(204)
	// 		return
	// 	}
	// 	c.Next()
	// })

	port := config.GetConfig().Common.Port
	if port == "" {
		port = "8000"
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
