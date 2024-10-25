package model

type DBResult struct {
	ID          int64
	IN_PROGRESS bool
	DATA        interface{}
}
