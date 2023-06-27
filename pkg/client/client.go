package client

type DbClient interface {
	Start() error
	Shutdown() error
	RunMigration() error
	GetConnection() (interface{}, error)
}
