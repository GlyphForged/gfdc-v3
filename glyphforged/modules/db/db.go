package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// ==================================================
//  DB Package
//
//  Handles opening, configuring, and seeding the db
//  -------------------------------------------------
//  Open(path string)
//  - Opens the DB. Creates it if it doesn't exist
//  params:
//  - path (string) - path to the db
//  returns:
//  - tuple consisting of a pointer to the db and err
//  -------------------------------------------------
//  Migrate()
//  - Ensures that the tables we want exist, creates
//  them if they don't
//  returns:
//  - error if present
//  -------------------------------------------------
//  Seed()
//  - Seeds the db with placeholder data
//  returns:
//  - error if one is generated
// ==================================================

// Opens or creates the sqlite DB file at path
func Open(path string) (*sql.DB, error) {
	// ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db dir err: %w", err)
	}

	// open Connection
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("Open db err: %w", err)
	}

	// set pool config
	db.SetMaxOpenConns(1) // single writer
	db.SetMaxIdleConns(1)

	// apply pragmas
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			db.Close()
			return nil, fmt.Errorf("apply pragma error %s: %s", p, err)
		}
	}

	// close if we get an error
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db err: %w", err)
	}

	return db, nil
}

// Migrate ensures tables/indexes exist
func Migrate(db *sql.DB) error {
	// create the schema
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY,
		slug TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		summary TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_projects_slug
		ON projects(slug);
	`

	// execute the schema
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("migrate schema err: %w", err)
	}

	return nil
}

// Seed inserts placeholder projects if the table is empty
func Seed(db *sql.DB) error {
	// figure out how many rows we have
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)

	// if none, seed the db
	if err == sql.ErrNoRows {
		// begin the transaction
		transaction, err := db.Begin()
		if err != nil {
			return err
		}

		// craft the statement with placeholders
		statement, err := transaction.Prepare(`
			INSERT INTO projects (slug, title, summary)
			VALUES (?, ?, ?)
		`)
		if err != nil {
			transaction.Rollback()
			return err
		}
		defer statement.Close()

		// set up the placeholders
		placeholders := []struct {
			slug, title, summary string
		}{
			{"project-alpha", "Project Alpha", "01001001 01101110 01110100 01100101 01101100 01101100 01101001 01100111 01100101 01101110 01100011 01100101 00100000 01101001 01110011 00100000 01110100 01101000 01100101 00100000 01101101 01100001 01101110 01101001 01100110 01100101 01110011 01110100 01100001 01110100 01101001 01101111 01101110 00100000 01101111 01100110 00100000 01110100 01101000 01100101 00100000 01001111 01101101 01101110 01101001 01110011 01110011 01101001 01100001 01101000"},
			{"project-delta", "Project Delta", "Omnis lorem ipsum dolor sit amet, consectetur mechanicus adipiscing elit. Sed do binary tempor incididunt ut labore et cogitator magna aliqua. Ut enim ad minim servos, quis nostrud exercitation standard template construct nisi ut aliquip ex ea commodo datasmith."},
			{"project-gamma", "Project Gamma", "Duis aute irure cog in reprehenderit in vox-caster velit esse cillum dolore eu fugiat nulla machine-spirit pariatur. Excepteur sint heretek non proident, sunt in culpa qui scrap-code mollit anim id est Mars."},
		}

		// execute the statement, 1 for each placeholder
		for _, p := range placeholders {
			if _, err := statement.Exec(p.slug, p.title, p.summary); err != nil {
				transaction.Rollback()
				return err
			}
		}

		// commit the shit
		return transaction.Commit()
	}

	// fallthrough if db is already seeded
	return nil
}
