package db

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	model "WebServer/internal/models/db/model"
)

func (w *Worker) GetUser(email string) (model.DBUser, error) {
	uri := w.url + "users/email/"
	type Request struct {
		Email string `json:"email"`
	}
	req := Request{
		Email: email,
	}
	body, err := json.Marshal(req)
	if err != nil {
		w.logger.Error("Error marshalling request to DB Users:", "err", err)
		return model.DBUser{}, err
	}
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(body))
	if err != nil {
		w.logger.Error("Error sendind request to DB Users:", "err", err)
		return model.DBUser{}, err
	}
	var user model.DBUser
	byteBody, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(byteBody, &user)
	if err != nil {
		w.logger.Error("Error unmarshalling response from DB Users:", "err", err)
		return model.DBUser{}, err
	}
	return user, nil
}

func (w *Worker) GetUserByID(id int) (model.DBUser, error) {
	uri := w.url + "users/id/"
	type Request struct {
		ID int `json:"id"`
	}
	req := Request{
		ID: id,
	}
	body, err := json.Marshal(req)
	if err != nil {
		w.logger.Error("Error marshalling request to DB Users:", "err", err)
		return model.DBUser{}, err
	}
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(body))
	if err != nil {
		w.logger.Error("Error sendind request to DB Users:", "err", err)
		return model.DBUser{}, err
	}
	var user model.DBUser
	byteBody, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(byteBody, &user)
	if err != nil {
		w.logger.Error("Error unmarshalling response from DB Users:", "err", err)
		return model.DBUser{}, err
	}
	return user, nil
}

func (w *Worker) CheckForRegistered(email string) bool {
	uri := w.url + "users/check/"
	type Request struct {
		Email string `json:"email"`
	}
	req := Request{
		Email: email,
	}
	body, err := json.Marshal(req)
	if err != nil {
		w.logger.Error("Error marshalling request to DB Users:", "err", err)
		return true
	}
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(body))
	if err != nil {
		w.logger.Error("Error sendind request to DB Users:", "err", err)
		return true
	}
	if resp.StatusCode != http.StatusOK {
		w.logger.Error("Status code is not 200:", "status", resp.Status)
	}

	type Response struct {
		Error  string `json:"error"`
		Status bool   `json:"status"`
	}

	byteBody, err := io.ReadAll(resp.Body)
	response := Response{}
	err = json.Unmarshal(byteBody, &response)
	if err != nil {
		w.logger.Error("Error unmarshalling response from DB Users:", "err", err)
		return true
	}
	if response.Error != "" {
		w.logger.Error("Error from DB Users:", "err", response.Error)
		return true
	}
	log.Println("Статус пользователя:", response.Status)
	return response.Status
}

func (w *Worker) Register(email, hashPassword, FirstName, LastName string) string {
	uri := w.url + "users/register"
	type Request struct {
		Email        string `json:"email"`
		HashPassword string `json:"password"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
	}
	type Response struct {
		Error  string `json:"error"`
		UserID int    `json:"id"`
	}
	request := Request{
		Email:        email,
		HashPassword: hashPassword,
		FirstName:    FirstName,
		LastName:     LastName,
	}
	byteBody, err := json.Marshal(request)
	if err != nil {
		w.logger.Error("Error marshalling request to DB Users:", "err", err)
		return ""
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(byteBody))
	if err != nil {
		w.logger.Error("Error sendind request to DB Users:", "err", err)
		return ""
	}
	byteBody, err = io.ReadAll(resp.Body)
	response := Response{}
	err = json.Unmarshal(byteBody, &response)
	if err != nil {
		w.logger.Error("Error unmarshalling response from DB Users:", "err", err)
		return ""
	}
	if response.Error != "" {
		w.logger.Error("Error from DB Users:", "err", response.Error)
		return ""
	}
	return strconv.Itoa(response.UserID)
}

func (w *Worker) CheckForLogin(email, hashPassword string) (status bool, user_id string) {
	uri := w.url + "users/compare"
	type Request struct {
		Email        string `json:"email"`
		HashPassword string `json:"password"`
	}
	type Response struct {
		Error  string `json:"error"`
		Status bool   `json:"status"`
		UserID int    `json:"id"`
	}
	request := Request{
		Email:        email,
		HashPassword: hashPassword,
	}
	byteBody, err := json.Marshal(request)
	if err != nil {
		w.logger.Error("Error marshalling request to DB Users:", "err", err)
		return false, ""
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(byteBody))
	if err != nil {
		w.logger.Error("Error sendind request to DB Users:", "err", err)
		return false, ""
	}
	byteBody, err = io.ReadAll(resp.Body)
	response := Response{}
	err = json.Unmarshal(byteBody, &response)
	if err != nil {
		w.logger.Error("Error unmarshalling response from DB Users:", "err", err)
		return false, ""
	}
	if response.Error != "" {
		w.logger.Error("Error from DB Users:", "err", response.Error)
		return false, ""
	}
	return response.Status, strconv.Itoa(response.UserID)
}

func (w *Worker) GetAllUsers() ([]model.DBUser, error) {
	uri := w.url + "users/all"
	resp, err := http.Get(uri)
	if err != nil {
		w.logger.Error("Error sendind request to DB Users:", "err", err)
		return nil, err
	}
	var users []model.DBUser
	byteBody, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(byteBody, &users)
	if err != nil {
		w.logger.Error("Error unmarshalling response from DB Users:", "err", err)
		return nil, err
	}
	return users, nil
}
