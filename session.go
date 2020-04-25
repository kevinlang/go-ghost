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

	response, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to establish session")
	}

	return nil
}
