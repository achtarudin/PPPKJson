package db_adapter

import (
	"context"
	"cutbray/pppk-json/internal/ports"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ ports.AdapterPort = &dbAdapter{}

type dbAdapter struct {
	dsn        string
	gorm       *gorm.DB
	gormConfig *gorm.Config
}

func (db *dbAdapter) Connect(ctx context.Context) error {
	database, err := gorm.Open(postgres.Open(db.dsn), db.gormConfig)

	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db.gorm = database
	return nil
}

func (db *dbAdapter) Disconnect(ctx context.Context) error {

	if db.gorm == nil {
		return nil
	}

	errChan := make(chan error, 1)
	go func() {
		sqlDB, err := db.gorm.DB()
		if err != nil {
			errChan <- err
			return
		}
		errChan <- sqlDB.Close()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}

}

func (g *dbAdapter) Value() any {
	return g.gorm
}

func (db *dbAdapter) IsReady() bool {
	if db.gorm == nil {
		return false
	}

	sqlDB, err := db.gorm.DB()
	if err != nil {
		return false
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		return false
	}

	return true
}

func New(dsn string, gormConfig *gorm.Config) *dbAdapter {
	return &dbAdapter{
		dsn:        dsn,
		gormConfig: gormConfig,
	}
}
