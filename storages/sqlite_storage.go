package storages

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/takahawk/shadownet/logger"
)

// DbDriverNameSqlite3 is name of SQLite3 driver
const DbDriverNameSqlite3 = "sqlite3"

var migrations = []string{
	`CREATE TABLE pipelines (
		name TEXT PRIMARY KEY,
		json TEXT NOT NULL
	)`,
}

type sqliteStorage struct {
	db     *sql.DB
	logger logger.Logger
}

// NewSqliteStorage returns new storage backed by SQLite database.
// It runs migrations after database connection, so that errors returned can
// be related to both DB connection and running migrations
func NewSqliteStorage(filename string, logger logger.Logger) (Storage, error) {
	db, err := sql.Open(DbDriverNameSqlite3, filename)

	if err != nil {
		logger.Errorf("%+v", err)
		return nil, err
	}

	storage := &sqliteStorage{
		db:     db,
		logger: logger,
	}

	err = storage.runMigrations()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

// ListPipelineJSONs returns all pipeline JSONs stored in SQLite database
func (ss *sqliteStorage) ListPipelineJSONs() ([]PipelineJSONsListEntry, error) {
	result := make([]PipelineJSONsListEntry, 0)
	rows, err := ss.db.Query("SELECT name, json FROM pipelines")
	if err != nil {
		ss.logger.Errorf("Error getting pipeline: %+v", err)
		return nil, err
	}

	for rows.Next() {
		var entry PipelineJSONsListEntry
		err = rows.Scan(&entry.Name, &entry.JSON)
		if err != nil {
			ss.logger.Errorf("Error getting pipeline: %+v", err)
			return nil, err
		}
		result = append(result, entry)
	}

	return result, nil
}

// SavePipelineJSON saves pipeline JSON in SQL database
func (ss *sqliteStorage) SavePipelineJSON(name string, json string) error {
	_, err := ss.db.Exec("INSERT INTO pipelines (name, json) VALUES (?, ?)", name, json)
	if err != nil {
		ss.logger.Errorf("Error saving pipeline: %+v", err)
		// mb more verbose logging?
		return err
	}
	return nil
}

// LoadPipelineJSON makes query to SQLite to get pipeline JSON
func (ss *sqliteStorage) LoadPipelineJSON(name string) (json string, err error) {
	row := ss.db.QueryRow("SELECT json FROM pipelines WHERE name = ?", name)
	var pipelineJson string
	err = row.Scan(&pipelineJson)
	if err != nil {
		ss.logger.Errorf("Error getting pipeline: %+v", err)
		return "", err
	}
	return pipelineJson, nil
}

// UpdatePipelineJSON makes query to SQLite to overwrite pipeline JSON with a
// given name
func (ss *sqliteStorage) UpdatePipelineJSON(name string, json string) error {
	// TODO: return error if pipeline is not exists
	_, err := ss.db.Exec("UPDATE pipelines SET json = ? WHERE name = ?", json, name)
	if err != nil {
		ss.logger.Errorf("Error updating pipeline: %+v", err)
		// mb more verbose logging?
		return err
	}
	return nil
}

// DeletePipelineJSON makes query to remove pipeline JSON with a given name
// from database
func (ss *sqliteStorage) DeletePipelineJSON(name string) error {
	_, err := ss.db.Exec("DELETE FROM pipelines WHERE name = ?", name)
	if err != nil {
		ss.logger.Errorf("Error deleting pipeline: %+v", err)
		// mb more verbose logging?
		return err
	}
	return nil
}

func (ss *sqliteStorage) getSchemaVersion() (int, error) {
	row := ss.db.QueryRow("PRAGMA schema_version")
	var version int
	err := row.Scan(&version)
	if err != nil {
		ss.logger.Errorf("%+v", err)
		return 0, err
	}
	return version, nil
}

func (ss *sqliteStorage) runMigrations() error {
	version, err := ss.getSchemaVersion()
	if err != nil {
		return err
	}

	if version == len(migrations) {
		ss.logger.Info("Database schema is up-to-date")
		return nil
	}

	ss.logger.Infof("Updating database schema from version %d to %d", version, len(migrations))

	for i := version; i < len(migrations); i++ {
		ss.logger.Infof("Applying migration %d: %s", i, migrations[i])
		_, err := ss.db.Exec(migrations[i])
		if err != nil {
			ss.logger.Errorf("%+v", err)
			return err
		}
	}

	return nil
}
