package ghost

import (
	"fmt"
	"time"
)

// AdminPostsService provides access to Post related functions in the Ghost Admin API.
type AdminPostsService adminService

// Role represents the role a user may have.
type Role struct {
	ID          *string    `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

// Author represents an author.
type Author struct {
	ID              *string    `json:"id"`
	Name            *string    `json:"name"`
	Slug            *string    `json:"slug"`
	Email           *string    `json:"email"`
	ProfileImage    *string    `json:"profile_image"`
	CoverImage      *string    `json:"cover_image"`
	Bio             *string    `json:"bio"`
	Website         *string    `json:"website"`
	Location        *string    `json:"location"`
	Facebook        *string    `json:"facebook"`
	Twitter         *string    `json:"twitter"`
	Accessibility   *string    `json:"accessibility"`
	Status          *string    `json:"status"`
	MetaTitle       *string    `json:"meta_title"`
	MetaDescription *string    `json:"meta_description"`
	Tour            *bool      `json:"tour"`
	LastSeen        *time.Time `json:"last_seen"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	Roles           []*Role    `json:"roles"`
	URL             *string    `json:"url"`
}

// Post represents a Ghost post.
type Post struct {
	Slug               *string    `json:"slug"`
	ID                 *string    `json:"id"`
	UUID               *string    `json:"uuid"`
	Title              *string    `json:"title"`
	Mobiledoc          *string    `json:"mobiledoc"`
	HTML               *string    `json:"html"`
	CommentID          *string    `json:"comment_id"`
	FeatureImage       *string    `json:"feature_image"`
	Featured           *bool      `json:"featured"`
	Status             *string    `json:"status"`
	Visibility         *string    `json:"visibility"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	PublishedAt        *time.Time `json:"published_at"`
	CustomExcerpt      *string    `json:"custom_excerpt"`
	CodeinjectionHead  *string    `json:"codeinjection_head"`
	CodeinjectionFoot  *string    `json:"codeinjection_foot"`
	CustomTemplate     *string    `json:"custom_template"`
	CanonicalURL       *string    `json:"canonical_url"`
	Tags               []*Tag     `json:"tags"`
	Authors            []*Author  `json:"authors"`
	PrimaryAuthor      *Author    `json:"primary_author"`
	PrimaryTag         *Tag       `json:"primary_tag"`
	URL                *string    `json:"url"`
	Excerpt            *string    `json:"excerpt"`
	ReadingTime        *int       `json:"reading_time"`
	OgImage            *string    `json:"og_image"`
	OgTitle            *string    `json:"og_title"`
	OgDescription      *string    `json:"og_description"`
	TwitterImage       *string    `json:"twitter_image"`
	TwitterTitle       *string    `json:"twitter_title"`
	TwitterDescription *string    `json:"twitter_description"`
	MetaTitle          *string    `json:"meta_title"`
	MetaDescription    *string    `json:"meta_description"`
}

func (p Post) String() string {
	return Stringify(p)
}

// PostsResponse is the structure of the Post response.
type PostsResponse struct {
	Posts []*Post
	Meta  *Meta
}

func (pr PostsResponse) String() string {
	return Stringify(pr)
}

// Get fetches a post by id.
func (s *AdminPostsService) Get(id string) (*Post, error) {
	u := fmt.Sprintf("posts/%v", id)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	postsResponse := new(PostsResponse)
	_, err = s.client.Do(req, postsResponse)
	if err != nil {
		return nil, err
	}

	return postsResponse.Posts[0], nil
}

// List fetches all posts via the ListParams.
func (s *AdminPostsService) List(listParams *ListParams) (*PostsResponse, error) {
	u, err := addOptions("posts/", listParams)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	postsResponse := new(PostsResponse)
	_, err = s.client.Do(req, postsResponse)
	if err != nil {
		return nil, err
	}

	return postsResponse, nil
}
