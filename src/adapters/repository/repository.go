package repository

func (s *repository) TestDb() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}
