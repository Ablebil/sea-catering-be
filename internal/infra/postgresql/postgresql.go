package postgresql

import (
	"errors"
	"log"

	conf "github.com/Ablebil/sea-catering-be/config"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PgErrorCode string

const (
	ErrUniqueViolation      PgErrorCode = "23505"
	ErrForeignKeyViolation  PgErrorCode = "23503"
	ErrNotNullViolation     PgErrorCode = "23502"
	ErrCheckViolation       PgErrorCode = "23514"
	ErrDeadlockDetected     PgErrorCode = "40P01"
	ErrSerializationFailure PgErrorCode = "40001"
)

func New(dsn string, conf *conf.Config) (*gorm.DB, error) {
	var LogLevel logger.LogLevel
	if conf.AppEnv == "production" {
		LogLevel = logger.Warn
	} else {
		LogLevel = logger.Info
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(LogLevel),
	})
}

func CheckError(err error, target any) bool {
	if err == nil {
		return false
	}

	switch t := target.(type) {
	case PgErrorCode:
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == string(t) {
			log.Printf("PostgreSQL Error [%s]: %s\n", pgErr.Code, pgErr.Message)
		}
	case error:
		if errors.Is(err, t) {
			log.Println("Error: ", err.Error())
			return true
		}
	}

	return false
}
