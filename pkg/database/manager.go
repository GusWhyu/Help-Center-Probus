package database

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/rafialfauzan/layrd/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBManager struct {
	mu  sync.RWMutex
	dbs map[string]*gorm.DB
}

func NewDBManager() *DBManager {
	return &DBManager{
		dbs: make(map[string]*gorm.DB),
	}
}

// Connect initializes a configured database connection and sets up connection pooling.
func (m *DBManager) Connect(dbCfg config.DatabaseConfig, isDebug bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbCfg.Host, dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Port, dbCfg.SSLMode,
	)

	logLevel := logger.Warn
	if isDebug {
		logLevel = logger.Info
	}

	gormCfg := &gorm.Config{
		Logger:         logger.Default.LogMode(logLevel),
		TranslateError: true,
	}

	dialector, err := databaseDialector(dbCfg, dsn)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open database [%s]: %w", dbCfg.Name, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB for [%s]: %w", dbCfg.Name, err)
	}

	// Set connection pooling parameters
	if dbCfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	}
	if dbCfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	}
	if dbCfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.ConnMaxLifetime) * time.Minute)
	}

	if err := sqlDB.Ping(); err != nil {
		if closeErr := sqlDB.Close(); closeErr != nil {
			log.Printf("[DBManager] Failed to close database [%s] after ping failure: %v", dbCfg.Name, closeErr)
		}
		return nil, fmt.Errorf("failed to ping database [%s]: %w", dbCfg.Name, err)
	}

	m.mu.Lock()
	m.dbs[dbCfg.Name] = db
	m.mu.Unlock()

	log.Printf("[DBManager] Connected to database [%s] using driver [%s]\n", dbCfg.Name, dbCfg.Driver)

	return db, nil
}

func databaseDialector(dbCfg config.DatabaseConfig, postgresDSN string) (gorm.Dialector, error) {
	switch dbCfg.Driver {
	case "", "postgres":
		return postgres.Open(postgresDSN), nil
	case "mysql":
		dsnConfig := mysqldriver.NewConfig()
		dsnConfig.User = dbCfg.User
		dsnConfig.Passwd = dbCfg.Password
		dsnConfig.Net = "tcp"
		dsnConfig.Addr = net.JoinHostPort(dbCfg.Host, dbCfg.Port)
		dsnConfig.DBName = dbCfg.DBName
		dsnConfig.Params = map[string]string{"charset": "utf8mb4"}
		dsnConfig.ParseTime = true
		return mysql.Open(dsnConfig.FormatDSN()), nil
	default:
		return nil, fmt.Errorf(
			"unsupported database driver [%s] for database [%s]; supported drivers are postgres and mysql",
			dbCfg.Driver,
			dbCfg.Name,
		)
	}
}

// Get returns the gorm.DB instance for a specific database name
func (m *DBManager) Get(name string) (*gorm.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	db, exists := m.dbs[name]
	if !exists {
		return nil, fmt.Errorf("database connection [%s] not found", name)
	}
	return db, nil
}

// MustGet returns gorm.DB or panics if not found
func (m *DBManager) MustGet(name string) *gorm.DB {
	db, err := m.Get(name)
	if err != nil {
		log.Panicf("[DBManager] Fatal: %v", err)
	}
	return db
}

// Default returns the default/primary database instance
func (m *DBManager) Default() *gorm.DB {
	return m.MustGet("default")
}

// CloseAll closes all active database connections
func (m *DBManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, db := range m.dbs {
		if sqlDB, err := db.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("[DBManager] Error closing database [%s]: %v\n", name, err)
			} else {
				log.Printf("[DBManager] Closed database [%s]\n", name)
			}
		}
	}
}
