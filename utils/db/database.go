package db

type Database interface {
	Connect() error
	Disconnect() error
	GetConnection() interface{}
}
