package model

import "time"

type DBResult struct {
	ID             int64     `json:"id"`
	OPERATION_ID   string    `json:"operationID"`
	IN_PROGRESS    bool      `json:"inProgress"`
	DATA           []byte    `json:"data"`
	OPERATION_TYPE string    `json:"type"`
	CREATION_DATE  time.Time `json:"creationDate"`
	FINISH_DATE    time.Time `json:"finishDate"`
	VERSION        int64     `json:"version"`
	USER_ID        int       `json:"userID"`
	FIRST_NAME     string    `json:"firstName"`
	LAST_NAME      string    `json:"lastName"`
	EMAIL          string    `json:"email"`
	USER_STATUS    int       `json:"userStatus"`
}

type DBUser struct {
	ID          int    `json:"id"`
	FIRSTNAME   string `json:"firstName"`
	LASTNAME    string `json:"lastName"`
	EMAIL       string `json:"email"`
	USER_STATUS int    `json:"userStatus"`
}
