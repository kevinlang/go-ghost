package ghost

type AdminDatabaseService adminService

type DatabaseResponse struct {
	DB map[string]interface{}
}

func (s *AdminDatabaseService) Get() (map[string]interface{}, error) {
	req, err := s.client.NewRequest("GET", "db", nil)
	if err != nil {
		return nil, err
	}

	dbResponse := new(DatabaseResponse)
	_, err = s.client.Do(req, dbResponse)
	if err != nil {
		return nil, err
	}

	return dbResponse.DB, nil
}
