package storages

import (
	"database/sql"
	"encoding/json"
	"errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/takahawk/shadownet/logger"
	"github.com/takahawk/shadownet/models"
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

// ListPipelineSpecs returns all pipeline specifications stored in SQLite database
func (ss *sqliteStorage) ListPipelineSpecs() ([]*models.PipelineSpec, error) {
	result := make([]*models.PipelineSpec, 0)
	rows, err := ss.db.Query("SELECT name, json FROM pipelines")
	if err != nil {
		ss.logger.Errorf("Error getting pipeline: %+v", err)
		return nil, err
	}

	for rows.Next() {
		var name, pipelineJson string
		err = rows.Scan(&name, &pipelineJson)
		if err != nil {
			ss.logger.Errorf("Error getting pipeline: %+v", err)
			return nil, err
		}

		var spec models.PipelineSpec
		err = json.Unmarshal([]byte(pipelineJson), &spec)
		if err != nil {
			ss.logger.Errorf("Error unmarshaling pipeline JSON: %+v", err)
			return nil, err
		}
		result = append(result, &spec)
	}

	return result, nil
}

// SavePipelineSpec saves pipeline specifications in SQL database
func (ss *sqliteStorage) SavePipelineSpec(spec *models.PipelineSpec) error {
	pipelineJSON, err := json.Marshal(spec)
	if err != nil {
		ss.logger.Errorf("Error marshaling pipeline to JSON: %+v", err)
		// mb more verbose logging?
		return err
	}
	if spec.Name == "" {
		ss.logger.Error("Empty name of pipeline specification")
		return errors.New("empty name of pipeline specificafion")
	}
	_, err = ss.db.Exec("INSERT INTO pipelines (name, json) VALUES (?, ?)", spec.Name, pipelineJSON)
	if err != nil {
		ss.logger.Errorf("Error saving pipeline: %+v", err)
		// mb more verbose logging?
		return err
	}
	return nil
}

// LoadPipelineSpec makes query to SQLite to get pipeline specification
func (ss *sqliteStorage) LoadPipelineSpec(name string) (*models.PipelineSpec, error) {
	row := ss.db.QueryRow("SELECT json FROM pipelines WHERE name = ?", name)
	var pipelineJson string
	err := row.Scan(&pipelineJson)
	if err != nil {
		ss.logger.Errorf("Error getting pipeline: %+v", err)
		return nil, err
	}

	var pipelineSpec models.PipelineSpec
	err = json.Unmarshal([]byte(pipelineJson), &pipelineSpec)
	if err != nil {
		ss.logger.Errorf("Error unmarshaling pipeline from JSON: %+v", err)
		return nil, err
	}
	return &pipelineSpec, nil
}

// UpdatePipelineSpec makes query to SQLite to overwrite pipeline
// specification with a given name
func (ss *sqliteStorage) UpdatePipelineSpec(spec *models.PipelineSpec) error {
	pipelineJSON, err := json.Marshal(spec)
	if err != nil {
		ss.logger.Errorf("Error marshaling pipeline to JSON: %+v", err)
		// mb more verbose logging?
		return err
	}
	// TODO: return error if pipeline is not exists
	_, err = ss.db.Exec("UPDATE pipelines SET json = ? WHERE name = ?", spec.Name, pipelineJSON)
	if err != nil {
		ss.logger.Errorf("Error updating pipeline: %+v", err)
		// mb more verbose logging?
		return err
	}
	return nil
}

// DeletePipelineJSON makes query to remove pipeline specification with a given
// name from database
func (ss *sqliteStorage) DeletePipelineSpec(name string) error {
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
