package model

type DBResult struct {
	ID           int    `db:"id"`
	OPERATION_ID string `db:"operation_id"`
	IN_PROGRESS  bool   `db:"in_progress"`
	DATA         []byte `db:"data"`
	Error        string `db:"error"`
}
