package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type PG struct {
	Host     string
	Username string
	Password string
	Database string
}

func (p *PG) Connect(ctx context.Context) (*pgx.Conn, error) {
	pgUrl := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", p.Username, p.Password, p.Host, p.Database)
	log.Println(p.Username)

	conn, err := pgx.Connect(ctx, pgUrl)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
