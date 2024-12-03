package model

import "time"

type DBResult struct {
	ID             int64     `db:"id" json:"id"`
	OPERATION_ID   string    `db:"operation_id" json:"operation_id"`
	IN_PROGRESS    bool      `db:"in_progress" json:"in_progress"`
	DATA           []byte    `db:"data" json:"data"`
	OPERATION_TYPE string    `db:"type" json:"type"`
	CREATION_DATE  time.Time `db:"creation_date" json:"creation_date"`
	FINISH_DATE    time.Time `db:"finish_date" json:"finish_date"`
	VERSION        int64     `db:"version" json:"version"`
	USER_ID        int       `db:"user_id" json:"user_id"`
	FIRST_NAME     string    `db:"first_name" json:"first_name"`
	LAST_NAME      string    `db:"last_name" json:"last_name"`
	EMAIL          string    `db:"email" json:"email"`
}

type DBUser struct {
	ID        int    `db:"id" json:"id"`
	FIRSTNAME string `db:"first_name" json:"firstName"`
	LASTNAME  string `db:"last_name" json:"lastName"`
	EMAIL     string `db:"email" json:"email"`
}
