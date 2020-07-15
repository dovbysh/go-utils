package migrations

import (
	"fmt"
	"github.com/dovbysh/tests_common/v3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var m PgMigration

func TestPgMigration(t *testing.T) {
	wd, _ := os.Getwd()
	var wg sync.WaitGroup
	var pgCloser func()
	o, pgCloser, _, _ := tests_common.PostgreSQLContainer(&wg)
	defer pgCloser()
	wg.Wait()

	p := strings.Split(o.Addr, ":")
	if len(p) != 2 {
		t.Fatalf("can not parse db addr: %s, parsed: %#v", o.Addr, p)
	}
	port, e := strconv.Atoi(p[1])
	if e != nil {
		t.Fatalf("can not parse db port, addr: %s, parsed: %#v, error: %s", o.Addr, p, e)
	}
	m = PgMigration{
		User:       o.User,
		Password:   o.Password,
		Database:   o.Database,
		Host:       p[0],
		Port:       port,
		TableName:  "_migrations",
		SchemaName: "",
	}
	m.SetMigrationDir(fmt.Sprintf("%s/postgres_test_migrations", wd))
	t.Run("Up", testUp)
	t.Run("Down", testDown)
	t.Run("New", testNew)
}

func testUp(t *testing.T) {
	s, e := m.Up()
	if e != nil {
		t.Fatal(e)
	}
	t.Log(s)
	assert.NotEmpty(t, s)
}

func testDown(t *testing.T) {
	s, e := m.Down()
	if e != nil {
		t.Fatal(e)
	}
	t.Log(s)
	assert.NotEmpty(t, s)
}

func testNew(t *testing.T) {
	s, e := m.NewMigration("zzz")
	if e != nil {
		t.Fatal(e)
	}
	t.Log(s)
	assert.NotEmpty(t, s)
	assert.FileExists(t, s)
	b, err := ioutil.ReadFile(s)
	assert.NoError(t, err)
	t.Log(string(b))
	assert.NoError(t, os.Remove(s))
}
