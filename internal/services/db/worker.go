package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	model "WebServer/internal/models/db/model"

	config "WebServer/internal/config/services/db"
)

type Worker struct {
	logger *slog.Logger
	url    string
}

func New(logger *slog.Logger) *Worker {
	return &Worker{
		logger: logger,
		url:    config.URL,
	}
}

func (w *Worker) RegisterOperation(uniqID string, operation_type string, user_id int) error {
	uri := w.url

	type Request struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		UserID int    `json:"user_id"`
	}

	var request Request
	request.ID = uniqID
	request.Type = operation_type
	request.UserID = user_id

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	response, err := http.Post(uri, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Error")
	}
	return nil
}

func (w *Worker) SetResult(uniqID string, data []byte) error {
	uri := w.url + "operation/" + uniqID
	response, err := http.Post(uri, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("Error")
	}
	return nil
}

func (w *Worker) GetResult(uniqID string) (dbResult model.DBResult, err error) {
	uri := w.url + "operation/" + uniqID
	response, err := http.Get(uri)
	if err != nil {
		return dbResult, err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&dbResult)
	return dbResult, err
}

func (w *Worker) GetAllOperations(limit string, operation_type string, operation_id string) (dbResult []model.DBResult, err error) {
	uri := w.url
	if limit == "" {
		limit = "0"
	}

	uri += "?limit=" + limit
	if operation_type != "" {
		uri += "&type=" + operation_type
	}
	if operation_id != "" {
		uri += "&id=" + operation_id
	}
	response, err := http.Get(uri)
	if err != nil {
		return dbResult, err
	}
	defer response.Body.Close()

	type Response struct {
		Data  []model.DBResult `json:"operations"`
		Error string           `json:"error"`
	}

	var responseData Response
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err == nil && responseData.Error != "" {
		err = fmt.Errorf(responseData.Error)
	}
	return responseData.Data, err
}

func (w *Worker) GetUserOperations(user_id int, limit int, operation_type string) (dbResult []model.DBResult, err error) {
	uri := w.url + "?userid=" + strconv.Itoa(user_id) + "&limit=" + strconv.Itoa(limit) + "&type=" + operation_type
	response, err := http.Get(uri)
	if err != nil {
		return dbResult, err
	}
	defer response.Body.Close()
	type Response struct {
		Data  []model.DBResult `json:"operations"`
		Error string           `json:"error"`
	}
	var responseData Response
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err == nil && responseData.Error != "" {
		err = fmt.Errorf(responseData.Error)
	}
	return responseData.Data, err
}

func (w *Worker) GetOperationID() (id string, err error) {
	uri := w.url + "id/"
	response, err := http.Get(uri)
	if err != nil {
		return id, err
	}
	defer response.Body.Close()
	type Response struct {
		ID    string `json:"id"`
		Error string `json:"error"`
	}
	var idR Response
	err = json.NewDecoder(response.Body).Decode(&idR)
	if err == nil && idR.Error != "" {
		err = fmt.Errorf(idR.Error)
	}
	return idR.ID, err
}

func (w *Worker) GetVersion(uniqID string) (version int64, err error) {
	uri := w.url + "operation/version/" + uniqID
	type Response struct {
		Version int64  `json:"version"`
		Error   string `json:"error"`
	}
	var versionR Response
	response, err := http.Get(uri)
	if err != nil {
		return version, err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&versionR)
	if err == nil && versionR.Error != "" {
		err = fmt.Errorf(versionR.Error)
	}
	return versionR.Version, err
}
