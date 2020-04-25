package ghost

import "time"

// Tag represents a post/page tag.
type Tag struct {
	ID              *string    `json:"id"`
	Name            *string    `json:"name"`
	Slug            *string    `json:"slug"`
	Description     *string    `json:"description"`
	FeatureImage    *string    `json:"feature_image"`
	Visibility      *string    `json:"visibility"`
	MetaTitle       *string    `json:"meta_title"`
	MetaDescription *string    `json:"meta_description"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	URL             *string    `json:"url"`
}

func (t Tag) String() string {
	return Stringify(t)
}
