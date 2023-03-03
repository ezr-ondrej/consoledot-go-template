package main

import (
	"consoledot-go-template/internal/config"
	"consoledot-go-template/internal/db"
	"consoledot-go-template/internal/logging"
	"context"

	"github.com/rs/zerolog/log"
)

//func init() {
//	random.SeedGlobal()
//}

func main() {
	ctx := context.Background()
	config.Initialize("config/api.env", "config/migrate.env")

	// initialize stdout logging and AWS clients first (cloudwatch is not available in init containers)
	logger, closeFn := logging.InitializeLogger()
	defer closeFn()
	log.Logger = logger

	err := db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	err = db.Migrate(ctx, "public")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error running migration")
		return
	}
}
