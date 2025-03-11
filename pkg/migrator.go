package pkg

import (
	"database/sql"
	"fmt"
	"interview/adapter"

	migrate "github.com/rubenv/sql-migrate"
)

type MigratorConfig struct {
	Podro       string `koanf:"podro"`
	MigrationDB string `koanf:"migration_db"`
	Dialect     string `koanf:"dialect"`
}

type Migrator struct {
	config          MigratorConfig
	podroMigrations *migrate.FileMigrationSource
}

func NewMigrator(config MigratorConfig) Migrator {
	migrations := &migrate.FileMigrationSource{
		Dir: config.Podro,
	}

	return Migrator{config: config, podroMigrations: migrations}
}

func (m Migrator) Up(SQLDBConfig adapter.SQLDBConfig) error {
	migrate.SetTable(m.config.MigrationDB)

	db, err := sql.Open(m.config.Dialect, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		SQLDBConfig.User, SQLDBConfig.Password, SQLDBConfig.Host, SQLDBConfig.Port, SQLDBConfig.DB))
	if err != nil {
		return err
	}

	n, err := migrate.Exec(db, m.config.Dialect, m.podroMigrations, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}

func (m Migrator) Down(SQLDBConfig adapter.SQLDBConfig) error {
	migrate.SetTable(m.config.MigrationDB)

	db, err := sql.Open(m.config.Dialect, fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		SQLDBConfig.User, SQLDBConfig.Password, SQLDBConfig.Host, SQLDBConfig.Port, SQLDBConfig.DB))
	if err != nil {
		return err
	}

	n, err := migrate.Exec(db, m.config.Dialect, m.podroMigrations, migrate.Down)
	if err != nil {
		return err
	}
	fmt.Printf("Rollbacked %d migrations!\n", n)
	return nil
}
