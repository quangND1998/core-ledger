package database

import (
	config "core-ledger/configs"
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string // disable, require, verify-ca, verify-full
	TimeZone string
}

func connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

var once sync.Once
var instance *gorm.DB

func Instance() *gorm.DB {
	once.Do(func() {
		host := config.Reader().Get("PG_HOST")
		user := config.Reader().Get("PG_USER")
		password := config.Reader().Get("PG_PASSWORD")
		dbName := config.Reader().Get("PG_DB")
		port := config.Reader().Get("PG_PORT")
		sslMode := config.Reader().Get("PG_SSLMODE")

		// Validate required fields
		if host == "" {
			panic("PG_HOST environment variable is required but not set")
		}
		if user == "" {
			panic("PG_USER environment variable is required but not set")
		}
		if password == "" {
			panic("PG_PASSWORD environment variable is required but not set")
		}
		if dbName == "" {
			panic("PG_DB environment variable is required but not set")
		}
		if port == "" {
			port = "5432" // default PostgreSQL port
		}
		if sslMode == "" {
			sslMode = "require" // default for production
		}

		// Build DSN with URL encoding for special characters
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
			host, user, password, dbName, port, sslMode,
		)
		// Log DSN without password for debugging
		dsnLog := fmt.Sprintf(
			"host=%s user=%s password=*** dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
			host, user, dbName, port, sslMode,
		)
		fmt.Printf("🔍 DSN connecting to: %s\n", dsnLog)

		// Try connection with retry logic
		var db *gorm.DB
		var err error
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			if i > 0 {
				fmt.Printf("🔄 Retrying connection (attempt %d/%d)...\n", i+1, maxRetries)
				time.Sleep(time.Second * 2)
			}
			db, err = connect(dsn)
			if err == nil {
				break
			}
			fmt.Printf("❌ Connection attempt %d failed: %v\n", i+1, err)
		}

		if err != nil {
			// Try with unencoded password (in case URL encoding breaks it)
			fmt.Println("⚠️  Trying with unencoded password...")
			dsnFallback := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
				host, user, password, dbName, port, sslMode,
			)
			db, err = connect(dsnFallback)
			if err != nil {
				panic(fmt.Errorf("failed to connect database after %d attempts: %w\nPlease check:\n- Database credentials (PG_USER, PG_PASSWORD) - password may contain special characters\n- Database name (PG_DB)\n- Network connectivity\n- Firewall rules\n- SSL certificate (if using SSL)", maxRetries, err))
			}
		}

		// Configure connection pool
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(100)
			sqlDB.SetConnMaxLifetime(time.Hour)
		}

		if config.GetConfig().Common.Mode != "production" {
			db = db.Debug()
		}

		instance = db
		fmt.Println("✅ Database connected successfully")
	})
	return instance
}

func WithConfig(cfg *PostgresConfig) *gorm.DB {
	if cfg.SSLMode == "" {
		cfg.SSLMode = "require"
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode,
	)

	db, err := connect(dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}
	return db
}
func SwitchSchema(schema string) *gorm.DB {
	if schema == "" {
		schema = "public"
	}
	db := Instance().Session(&gorm.Session{}) // clone session mới
	db.Exec("SET search_path TO " + schema)
	return db
}
