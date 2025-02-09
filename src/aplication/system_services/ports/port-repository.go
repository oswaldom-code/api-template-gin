package ports

type Repository interface {
	TestDb() error
}
