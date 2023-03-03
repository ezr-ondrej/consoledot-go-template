package db

import (
	"consoledot-go-template/internal/config"
	"context"
	"fmt"
	"net/url"

	pgxlog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/georgysavva/scany/v2"
)

// Pool is the main connection pool for the whole application
var Pool *pgxpool.Pool

func Initialize(ctx context.Context, schema string) error {
	if schema == "" {
		schema = "public"
	}

	// register and setup logging configuration
	connStr := getConnString("postgres", schema)
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("unable to parse db configuration: %w", err)
	}

	logLevel, configErr := tracelog.LogLevelFromString(config.Logging.DatabaseLevel)
	if configErr != nil {
		return fmt.Errorf("cannot parse db log level configuration: %w", configErr)
	}

	if logLevel > 0 {
		zeroLogger := pgxlog.NewLogger(log.Logger,
			pgxlog.WithContextFunc(func(ctx context.Context, logWith zerolog.Context) zerolog.Context {
				//traceId := ctxval.TraceId(ctx)
				//if traceId != "" {
				//	logWith = logWith.Str("trace_id", traceId)
				//}
				//accountId := ctxval.AccountIdOrNil(ctx)
				//if accountId != 0 {
				//	logWith = logWith.Int64("account_id", accountId)
				//}
				return logWith
			}))
		poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   zeroLogger,
			LogLevel: logLevel,
		}
	}

	Pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = Pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("unable to ping the database: %w", err)
	}
	return nil
}

func Close() {
	log.Logger.Info().Msg("Closing all database connections")
	Pool.Close()
}

func getConnString(prefix, schema string) string {
	if len(config.Database.Password) > 0 {
		return fmt.Sprintf("%s://%s:%s@%s:%d/%s?search_path=%s",
			prefix,
			url.QueryEscape(config.Database.User),
			url.QueryEscape(config.Database.Password),
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
			schema)
	} else {
		return fmt.Sprintf("%s://%s@%s:%d/%s?search_path=%s",
			prefix,
			url.QueryEscape(config.Database.User),
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
			schema)
	}
}
