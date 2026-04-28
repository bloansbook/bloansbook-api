package database

import (
	"context"
	"log"

	"github.com/bloansbook/bloansbook-api/pkg/config"
	"github.com/bloansbook/bloansbook-api/pkg/sysmsg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	var err error

	config, _ := pgxpool.ParseConfig(config.ApplicationConfig.Database.URL)

	config.ConnConfig.StatementCacheCapacity = 0
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("%s: %v\n", sysmsg.CannotConnect, err)
	}

	if err = Pool.Ping(context.Background()); err != nil {
		log.Fatalf("%s: %v\n", sysmsg.CannotPing, err)
	}

	log.Println(sysmsg.ConnectionSuccessful)
}

func Close() {
	Pool.Close()
	log.Println(sysmsg.ConnectionClosed)
}
