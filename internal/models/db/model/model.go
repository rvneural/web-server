package model

import "time"

type DBResult struct {
	ID             int64     `db:"id"`
	OPERATION_ID   string    `db:"operation_id"`
	IN_PROGRESS    bool      `db:"in_progress"`
	DATA           []byte    `db:"data"`
	Error          string    `db:"error"`
	OPERATION_TYPE string    `db:"type"`
	CREATION_DATE  time.Time `db:"creation_date"`
	FINISH_DATE    time.Time `db:"finish_date"`
}
