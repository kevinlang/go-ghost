package ghost

import (
	"fmt"
	"net/http"
)

// AdminSessionService handles establishing a cookie-based session with Ghost.
type AdminSessionService adminService

// userCredentials are the credentials of the user used to establish a session.
type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Create creates the session. The cookie should be set in the underlying
// http.Client cookiejar, allowing use of the session for the duration of the client.
func (s *AdminSessionService) Create(username, password string) error {
	creds := &userCredentials{
		Username: username,
		Password: password,
	}
	req, err := s.client.NewRequest("POST", "session/", creds)
	if err != nil {
		return err
	}

	// we want to read the entire response stream or we may run into a
	// race condition and have our next call hit 403 because the session token
	// has not yet been persisted due to quirks of express-router.
	// same underlying cause as https://github.com/expressjs/session/issues/360
	var body interface{}
	response, err := s.client.Do(req, body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to establish session")
	}

	return nil
}
