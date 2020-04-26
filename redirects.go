package ghost

import (
	"encoding/json"
	"mime/multipart"
)

// AdminRedirectsService handles downloading and uploading the redirects.json file.
type AdminRedirectsService adminService

// Redirect is a single redirect entry.
type Redirect struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Download fetches the redirectsq
func (s *AdminRedirectsService) Download() ([]*Redirect, error) {
	req, err := s.client.NewRequest("GET", "redirects/json", nil)
	if err != nil {
		return nil, err
	}

	var redirects []*Redirect
	_, err = s.client.Do(req, &redirects)
	if err != nil {
		return nil, err
	}

	return redirects, nil
}

// Upload uploads the redirects.
func (s *AdminRedirectsService) Upload(redirects []*Redirect) error {
	redirectsWriter := func(mpw *multipart.Writer) error {
		part, err := createFormFile(mpw, "redirects", "redirects.json", "application/json")
		if err != nil {
			return err
		}
		enc := json.NewEncoder(part)
		enc.SetEscapeHTML(false)
		return enc.Encode(redirects)
	}

	req, err := s.client.NewUploadRequest("redirects/json", redirectsWriter, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(req, nil)
	return err
}
