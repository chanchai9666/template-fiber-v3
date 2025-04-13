package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config represents database configuration
type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Schema   string
}

// Service represents a service that interacts with a database.
type Service interface {
	Health() map[string]string
	Close() error
	DB() *gorm.DB
}

type service struct {
	db        *gorm.DB
	config    Config
	stopRetry chan struct{}
}

var dbInstances = make(map[string]*service)

func New(name string, config Config) Service {
	if db, exists := dbInstances[name]; exists {
		return db
	}

	if err := validateConfig(config); err != nil {
		log.Fatalf("Invalid DB config for %s: %v", name, err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		config.Host, config.Username, config.Password, config.Database, config.Port, config.Schema)

	gormConfig := &gorm.Config{
		Logger:          logger.Default.LogMode(logger.Info),
		TranslateError:  true, //ถ้าเป็น true, GORM จะทำการแปลข้อผิดพลาดจากฐานข้อมูลเป็นข้อผิดพลาดที่สามารถเข้าใจได้ใน Go.
		CreateBatchSize: 150,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database %s: %v", name, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from GORM: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	svc := &service{
		db:        db,
		config:    config,
		stopRetry: make(chan struct{}),
	}

	dbInstances[name] = svc
	return svc
}

func (s *service) Health() map[string]string {
	stats := make(map[string]string)

	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = "cannot get sql.DB"
		return stats
	}

	if err := sqlDB.Ping(); err != nil {
		stats["status"] = "down"
		stats["error"] = err.Error()
		go s.retryConnection()
		return stats
	}

	dbStats := sqlDB.Stats()
	stats["status"] = "up"
	stats["open_connections"] = fmt.Sprintf("%d", dbStats.OpenConnections)
	stats["in_use"] = fmt.Sprintf("%d", dbStats.InUse)
	stats["idle"] = fmt.Sprintf("%d", dbStats.Idle)
	stats["max_open_connections"] = fmt.Sprintf("%d", dbStats.MaxOpenConnections)

	return stats
}

func (s *service) Close() error {
	close(s.stopRetry)

	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("cannot close db: %v", err)
	}
	log.Printf("Disconnected from DB: %s", s.config.Database)
	return sqlDB.Close()
}

func (s *service) retryConnection() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.connect(); err != nil {
				log.Printf("Retry DB connection failed: %v", err)
			} else {
				log.Println("Successfully reconnected to the database")
				return
			}
		case <-s.stopRetry:
			return
		}
	}
}

func (s *service) connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		s.config.Host, s.config.Username, s.config.Password, s.config.Database, s.config.Port, s.config.Schema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return err
	}

	s.db = db
	return nil
}

func (s *service) DB() *gorm.DB {
	return s.db
}

func validateConfig(config Config) error {
	required := map[string]string{
		"Host":     config.Host,
		"Port":     config.Port,
		"Database": config.Database,
		"Username": config.Username,
		"Password": config.Password,
		"Schema":   config.Schema,
	}

	var missing []string
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}

	return nil
}
