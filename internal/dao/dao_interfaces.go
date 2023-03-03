package dao

import (
	"consoledot-go-template/internal/models"
	"context"
)

var GetHelloDao func(ctx context.Context) HelloDao

// HelloDao groups access methods for access to state of hello.
type HelloDao interface {
	RecordHello(ctx context.Context, message *models.Hello) error
}
