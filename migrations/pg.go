package migrations

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

type PgMigration struct {
	User     string
	Password string
	Database string
	Host     string
	Port     int

	Direction       migrate.MigrationDirection
	Limit           int
	MigrationSource migrate.MigrationSource
	TableName       string
	SchemaName      string
}

func (o PgMigration) Apply() ([]string, error) {
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

	n, err := migrate.ExecMax(db, "postgres", o.MigrationSource, o.Direction, o.Limit)
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
