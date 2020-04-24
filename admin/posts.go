package ghostadmin

import (
	"context"
	"fmt"
)

type PostsService service

// Get fetches a post by id.
func (s *PostsService) Get(ctx context.Context, id string) (*Post, error) {
	u := fmt.Sprintf("posts/%v", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	organization := new(Organization)
	resp, err := s.client.Do(ctx, req, organization)
	if err != nil {
		return nil, resp, err
	}

	return organization, resp, nil
}
