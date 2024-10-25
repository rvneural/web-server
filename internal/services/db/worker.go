package db

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
