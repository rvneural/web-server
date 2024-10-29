package db

import (
	"fmt"

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

func NewWorker(host, port, login, password, db_name, table_name string) *Worker {
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

func (w *Worker) RegisterOperation(uniqID string) error {

	if len(uniqID) == 0 || len(uniqID) > 35 {
		return fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO $1 (operation_id, in_progress) VALUES ($2, $3)", w.table_name, uniqID, true)
	return err
}

func (w *Worker) SetResult(uniqID string, data []byte) error {

	if len(uniqID) == 0 || len(uniqID) > 35 {
		return fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE $1 SET result=$2, in_progress=$3 WHERE operation_id = $4", w.table_name, data, false, uniqID)
	return err
}

func (w *Worker) GetResult(uniqID string) (dbResult model.DBResult, err error) {

	if len(uniqID) == 0 || len(uniqID) > 35 {
		return model.DBResult{}, fmt.Errorf("uniqID is empty or too big")
	}

	db, err := w.connectToDB()
	if err != nil {
		return model.DBResult{}, err
	}
	defer db.Close()

	dbResults := []model.DBResult{}

	err = db.Select(&dbResults, "SELECT * FROM $1 WHERE operation_id = $2 LIMIT 1", w.table_name, uniqID)
	if err != nil {
		return model.DBResult{}, err
	}

	if len(dbResults) == 0 {
		return model.DBResult{}, fmt.Errorf("no results")
	}

	return dbResults[0], nil
}
