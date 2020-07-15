package migrations

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"os"
	"path"
	"strings"
	"time"
)

type Direction int

const (
	Up Direction = iota
	Down
)

var direction map[Direction]migrate.MigrationDirection

func init() {
	direction = map[Direction]migrate.MigrationDirection{
		Up:   migrate.Up,
		Down: migrate.Down,
	}
}

type PgMigration struct {
	User            string
	Password        string
	Database        string
	Host            string
	Port            int
	Limit           int
	TableName       string
	SchemaName      string
	migrationSource migrate.MigrationSource
	direction       Direction
	dir             string
}

func (o *PgMigration) SetMigrationDir(dir string) *PgMigration {
	o.dir = dir
	o.migrationSource = migrate.FileMigrationSource{Dir: dir}
	return o
}

func (o *PgMigration) Up() ([]string, error) {
	o.direction = Up
	return o.apply()
}

func (o *PgMigration) Down() ([]string, error) {
	o.direction = Down
	return o.apply()
}

func (o *PgMigration) apply() ([]string, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"dbname=%s host=%s port=%d user=%s password=%s sslmode=disable",
			o.Database,
			o.Host,
			o.Port,
			o.User,
			o.Password,
		),
	)
	if err != nil {
		return nil, err
	}

	migrate.SetTable(o.TableName)
	if o.SchemaName != "" {
		migrate.SetSchema(o.SchemaName)
	}

	n, err := migrate.ExecMax(db, "postgres", o.migrationSource, direction[o.direction], o.Limit)
	if err != nil {
		return nil, fmt.Errorf("Migration failed: %s", err)
	}

	res := make([]string, 0, 1)
	if n == 1 {
		res = append(res, "Applied 1 migration")
	} else {
		res = append(res, fmt.Sprintf("Applied %d migrations", n))
	}
	return res, nil
}

func (o *PgMigration) NewMigration(name string) (string, error) {
	templateContent := `
-- +migrate Up

-- +migrate Down
`

	if _, err := os.Stat(o.dir); os.IsNotExist(err) {
		return "", err
	}

	fileName := fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102150405"), strings.TrimSpace(name))
	pathName := path.Join(o.dir, fileName)
	f, err := os.Create(pathName)

	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.WriteString(templateContent); err != nil {
		return "", err
	}

	return pathName, nil
}
