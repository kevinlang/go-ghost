package ghost

import (
	"fmt"
)

// AdminDatabaseService handles fetching and uploading the database.
type AdminDatabaseService adminService

// DatabaseMeta is metadata of the source of the db dump.
type DatabaseMeta struct {
	Version    string `json:"version,omitempty"`
	ExportedOn int64  `json:"exported_on,omitempty"`
}

// DatabaseData is the actual data of the database.
//type DatabaseData map[string]interface{}

// Database is the representation of the database, with meta info.
type Database struct {
	Meta *DatabaseMeta          `json:"meta"`
	Data map[string]interface{} `json:"data"`
}

// databaseWrapper is the form of the response we get that we later flatten.
type databaseWrapper struct {
	DB []*Database `json:"db"`
}

// Get retrieves the database export from the instance.
func (s *AdminDatabaseService) Get() (*Database, error) {
	req, err := s.client.NewRequest("GET", "db", nil)
	if err != nil {
		return nil, err
	}

	dbWrapper := new(databaseWrapper)
	_, err = s.client.Do(req, dbWrapper)
	if err != nil {
		return nil, err
	}

	if len(dbWrapper.DB) != 1 {
		return nil, fmt.Errorf("received unexpected response format")
	}
	return dbWrapper.DB[0], nil
}
