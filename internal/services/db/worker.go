package db

import (
	"fmt"
	"strings"

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
}

func New(host, port, login, password, db_name, table_name string) *Worker {
	return &Worker{
		host:       host,
		port:       port,
		login:      login,
		password:   password,
		db_name:    db_name,
		table_name: table_name,
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
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO "+w.table_name+" (operation_id, in_progress, type) VALUES ($1, $2, $3)", uniqID, true, operation_type)
	return err
}

func (w *Worker) SetResult(uniqID string, data []byte) error {

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE "+w.table_name+" SET data = $1, in_progress = $2 WHERE operation_id = $3", data, false, uniqID)
	return err
}

func (w *Worker) GetResult(uniqID string) (dbResult model.DBResult, err error) {

	uniqID = strings.TrimSpace(uniqID)
	if len(uniqID) == 0 || len(uniqID) > 35 {
		return model.DBResult{}, fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		return model.DBResult{}, err
	}
	defer db.Close()

	dbResults := make([]model.DBResult, 0, 2)

	err = db.Select(&dbResults, "SELECT * FROM "+w.table_name+" WHERE operation_id = $1", uniqID)
	if err != nil {
		return model.DBResult{}, err
	}

	if len(dbResults) == 0 {
		return model.DBResult{}, fmt.Errorf("no results")
	}

	return dbResults[0], nil
}
