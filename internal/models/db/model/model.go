package model

type DBResult struct {
	OPERATION_ID int    `db:"operation_id"`
	IN_PROGRESS  bool   `db:"in_progress"`
	DATA         []byte `db:"data"`
}
