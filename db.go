package ghost

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
)

// AdminDatabaseService handles fetching and uploading the database.
type AdminDatabaseService adminService

// DatabaseMeta is metadata of the source of the db dump.
type DatabaseMeta struct {
	Version    string `json:"version,omitempty"`
	ExportedOn int64  `json:"exported_on,omitempty"`
}

// Database is the representation of the database, with meta info.
type Database struct {
	Meta *DatabaseMeta          `json:"meta"`
	Data map[string]interface{} `json:"data"`
}

// DatabaseImportProblem represents any issues or strageness encountered during import.
type DatabaseImportProblem struct {
	Message string                 `json:"message"`
	Help    string                 `json:"help"`
	Context string                 `json:"context"`
	Err     map[string]interface{} `json:"err"`
}

type databaseImportWrapper struct {
	Problems []*DatabaseImportProblem `json:"problems"`
}

// databaseWrapper is the form of the response we get that we later flatten.
type databaseWrapper struct {
	DB []*Database `json:"db"`
}

// Export the database.
func (s *AdminDatabaseService) Export() (*Database, error) {
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

// Import the database. Returns the list of problems (warnings), if any.
func (s *AdminDatabaseService) Import(db *Database) ([]*DatabaseImportProblem, error) {
	dbPartWriter := func(mpw *multipart.Writer) error {
		part, err := createFormFile(mpw, "importfile", "ghost.json", "application/json")
		if err != nil {
			return err
		}
		enc := json.NewEncoder(part)
		enc.SetEscapeHTML(false)
		return enc.Encode(db)
	}

	req, err := s.client.NewUploadRequest("db", dbPartWriter, nil)
	if err != nil {
		return nil, err
	}

	wrapper := new(databaseImportWrapper)
	_, err = s.client.Do(req, wrapper)
	if err != nil {
		return nil, err
	}

	return wrapper.Problems, nil
}
