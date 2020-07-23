package tlog

import (
	"github.com/dovbysh/tests_common/v3"
	"github.com/go-pg/pg/v9"
	"sync"
	"testing"
	"time"
)

type Book struct {
	tableName   struct{} `pg:"books"`
	ID          uint64   `pg:",pk"`
	CreatedAt   time.Time
	PublishedAt time.Time
	Name        string
	Description string
}

func TestExample(t *testing.T) {
	var wg sync.WaitGroup
	var pgCloser func()
	o, pgCloser, _, _ := tests_common.PostgreSQLContainer(&wg)
	defer pgCloser()

	wg.Wait()

	db := pg.Connect(&pg.Options{
		User:     o.User,
		Password: o.Password,
		Database: o.Database,
		Addr:     o.Addr,
	})
	db.AddQueryHook(NewShowQuery(t))
	db.CreateTable((*Book)(nil), nil)
}
