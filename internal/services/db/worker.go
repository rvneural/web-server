package db

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	model "WebServer/internal/models/db/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Worker struct {
	host       string
	port       string
	login      string
	password   string
	db_name    string
	table_name string
	logger     *slog.Logger
}

func New(host, port, login, password, db_name, table_name string, logger *slog.Logger) *Worker {
	return &Worker{
		host:       host,
		port:       port,
		login:      login,
		password:   password,
		db_name:    db_name,
		table_name: table_name,
		logger:     logger,
	}
}

func (w *Worker) connectToDB() (*sqlx.DB, error) {
	connectionData := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s port=%s", w.login, w.db_name, w.password, w.host, w.port)
	return sqlx.Connect("postgres", connectionData)
}

func (w *Worker) RegisterOperation(uniqID string, operation_type string) error {

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO "+w.table_name+" (operation_id, in_progress, type, creation_date, finish_date, version) VALUES ($1, $2, $3, $4, $5, $6)", uniqID, true, operation_type, time.Now(), time.Now(), 0)
	if err != nil {
		w.logger.Error("Insert operation to DataBase", "error", err)
	}
	return err
}

func (w *Worker) SetResult(uniqID string, data []byte) error {

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE "+w.table_name+" SET data = $1, in_progress = $2, finish_date = $3, version = $4 WHERE operation_id = $5", data, false, time.Now(), 1, uniqID)
	if err != nil {
		w.logger.Error("Update operation to DataBase", "error", err)
	}
	return err
}

func (w *Worker) GetResult(uniqID string) (dbResult model.DBResult, err error) {

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return model.DBResult{}, fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return model.DBResult{}, err
	}
	defer db.Close()

	dbResults := make([]model.DBResult, 0, 2)

	err = db.Select(&dbResults, "SELECT * FROM "+w.table_name+" WHERE operation_id = $1", uniqID)
	if err != nil {
		w.logger.Error("Get operation from DataBase", "error", err)
		return model.DBResult{}, err
	}

	if len(dbResults) == 0 {
		return model.DBResult{}, fmt.Errorf("no results")
	}

	return dbResults[0], nil
}

func (w *Worker) GetAllOperations(limit int, operation_type string) (dbResult []model.DBResult, err error) {
	dbResult = make([]model.DBResult, 0, 10)
	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return dbResult, err
	}
	defer db.Close()

	request := "SELECT * FROM " + w.table_name

	if operation_type != "" {
		request += " WHERE type = '" + strings.ToLower(strings.TrimSpace(operation_type)) + "'"
	}

	request += " ORDER BY id DESC"

	if limit > 0 {
		request += " LIMIT " + strconv.Itoa(limit)
	}

	err = db.Select(&dbResult, request)
	if err != nil {
		w.logger.Error("Get operation from DataBase", "error", err)
	}
	return dbResult, err
}

func (w *Worker) GetOperation(uniqID string) (dbResult model.DBResult, err error) {
	uniqID = strings.ToLower(strings.TrimSpace(uniqID))
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return model.DBResult{}, fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return model.DBResult{}, err
	}
	defer db.Close()

	err = db.Get(&dbResult, "SELECT * FROM "+w.table_name+" WHERE operation_id = $1 LIMIT 1", uniqID)
	if err != nil {
		w.logger.Error("Get operation from DataBase", "error", err)
	}
	return dbResult, err
}

func (w *Worker) UpdateResult(uniqID string, data []byte) (err error) {

	w.logger.Info("UpdateResult", "uniqID", uniqID, "data", string(data))
	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE "+w.table_name+" SET data=$1, version=version+1 WHERE operation_id = $2", data, uniqID)
	return err
}

func (w *Worker) GetVersion(uniqID string) (version int64, err error) {

	type versionLabel struct {
		Version int64 `db:"version"`
	}

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return 0, fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		w.logger.Error("Connection to DataBase", "error", err)
		return 0, err
	}

	var dbLabel = versionLabel{}

	defer db.Close()
	err = db.Get(&dbLabel, "SELECT version FROM "+w.table_name+" WHERE operation_id = $1 LIMIT 1", uniqID)
	return dbLabel.Version, err
}
