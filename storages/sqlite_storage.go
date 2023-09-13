package storages

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/takahawk/shadownet/logger"
)

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

func (ss *sqliteStorage) SavePipelineJSON(name string, json string) error {
	_, err := ss.db.Exec("INSERT INTO pipelines (name, json) VALUES (?, ?)", name, json)
	if err != nil {
		ss.logger.Errorf("Error saving pipeline: %+v", err)
		// mb more verbose logging?
		return err
	}
	return nil
}

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
