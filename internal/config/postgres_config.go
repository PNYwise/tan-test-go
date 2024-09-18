package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
)

var db *pgx.Conn

func getDatabaseConn(ctx context.Context, conf *viper.Viper) *pgx.Conn {

	dbConfig := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		conf.GetString("pg.username"),
		conf.GetString("pg.password"),
		conf.GetString("pg.host"),
		conf.GetInt("pg.port"),
		conf.GetString("pg.name"),
	)
	connConfig, err := pgx.ParseConfig(dbConfig)
	if err != nil {
		panic(err)
	}

	newDB, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		panic(err)
	}

	if err := newDB.Ping(ctx); err != nil {
		panic(err)
	}
	fmt.Printf("Connected to Database \n")
	db = newDB
	return db
}
func DbConn(ctx context.Context, conf *viper.Viper) *pgx.Conn {
	if db == nil {
		db = getDatabaseConn(ctx, conf)
	}
	return db
}
