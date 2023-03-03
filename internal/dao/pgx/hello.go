package pgx

import (
	"consoledot-go-template/internal/dao"
	"consoledot-go-template/internal/db"
	"consoledot-go-template/internal/models"
	"context"
	"fmt"
)

func init() {
	dao.GetHelloDao = getHelloDao
}

type helloDaoPgx struct{}

func getHelloDao(ctx context.Context) dao.HelloDao {
	return &helloDaoPgx{}
}

func (x *helloDaoPgx) RecordHello(ctx context.Context, hello *models.Hello) error {
	query := `
		INSERT INTO hellos (from, to, message)
		VALUES ($1, $2, $3) RETURNING id`

	err := db.Pool.QueryRow(ctx, query, hello.From, hello.To, hello.Message).Scan(&hello.ID)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	return nil
}
