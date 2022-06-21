package database

import (
	"context"
	"database/sql"
	"embed"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
	"sort"
	"time"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Database struct {
	database *sql.DB
}

func Open() (*Database, error) {
	// Create an in-memory DB, by switching the driver and DSN this can be changed to Postgres or another RDBMS
	db, err := sql.Open("sqlite3", "file:chat.db?cache=private&mode=memory")
	if err != nil {
		return nil, err
	}

	// Create the struct and run migrations
	database := &Database{database: db}
	err = database.migrate()

	return database, err
}

// For simplicity, all migrations are idempotent
func (db *Database) migrate() error {
	files, err := migrations.ReadDir("migrations")
	if err != nil {
		return err
	}

	// Sort the slice of files by the filename in ascending order
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})

	// Begin a transaction to run migrations in
	migrationContext, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	tx, err := db.database.BeginTx(migrationContext, nil)
	for _, file := range files {
		log.Printf("Running migration %s\n", file.Name())
		if file.IsDir() {
			continue
		}
		var migrationData []byte
		migrationData, err = migrations.ReadFile(filepath.Join("migrations", file.Name()))
		if err != nil {
			log.Printf("Error reading migration %s\n", err.Error())
			break
		}

		_, err = tx.ExecContext(migrationContext, string(migrationData))
		if err != nil {
			log.Printf("Error running migration %s\n", err.Error())
			break
		}
	}

	// Commit or rollback the migration transaction
	if err == nil {
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction %s\n", err.Error())
		}
	} else {
		log.Println("Error occurred while migrating, rolling back transaction")
		if err := tx.Rollback(); err != nil {
			log.Printf("Error rolling back transaction %s\n", err.Error())
		}
	}
	cancel()

	return err
}
