package tlog

import (
	"github.com/dovbysh/tests_common/v3"
	pg10 "github.com/go-pg/pg/v10"
	pg9 "github.com/go-pg/pg/v9"
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

	db9 := pg9.Connect(&pg9.Options{
		User:     o.User,
		Password: o.Password,
		Database: o.Database,
		Addr:     o.Addr,
	})
	db9.AddQueryHook(NewShowQuery(t))
	db9.CreateTable((*Book)(nil), nil)
	db9.DropTable((*Book)(nil), nil)

	db10 := pg10.Connect(&pg10.Options{
		User:     o.User,
		Password: o.Password,
		Database: o.Database,
		Addr:     o.Addr,
	})
	db10.AddQueryHook(NewShowQuery10(t))
	db10.CreateTable((*Book)(nil), nil)
}
