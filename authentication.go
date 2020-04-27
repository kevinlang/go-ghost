package ghost

import (
	"fmt"
	"net/http"
)

// AdminAuthenticationService handles setting up, invitations, and password resets.
type AdminAuthenticationService adminService

// SetupDetails is the information needed to setup the Ghost instance.
type SetupDetails struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	BlogTitle string `json:"blogTitle"`
}

type setupWrapper struct {
	Setup []*SetupDetails `json:"setup"`
}

// Setup initializes the Ghost instance.
func (s *AdminAuthenticationService) Setup(details *SetupDetails) error {
	wrapper := &setupWrapper{
		Setup: []*SetupDetails{details},
	}
	req, err := s.client.NewRequest("POST", "authentication/setup", wrapper)
	if err != nil {
		return err
	}

	response, err := s.client.Do(req, nil)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to setup")
	}
	return nil
}
