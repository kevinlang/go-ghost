package ghost

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPost_marshall(t *testing.T) {

	primaryTag := &Tag{
		ID:              String("5ddc9063c35e7700383b27e0"),
		Name:            String("Getting Started"),
		Slug:            String("getting-started"),
		Description:     String("a"),
		FeatureImage:    String("b"),
		Visibility:      String("public"),
		MetaTitle:       String("c"),
		MetaDescription: String("d"),
		CreatedAt:       Time("2019-11-26T02:39:31.000Z"),
		UpdatedAt:       Time("2019-11-26T02:39:31.000Z"),
		URL:             String("https://docs.ghost.io/tag/getting-started/"),
	}
	primaryAuthor := &Author{
		ID:              String("5951f5fca366002ebd5dbef7"),
		Name:            String("Ghost"),
		Slug:            String("ghost-user"),
		Email:           String("info@ghost.org"),
		ProfileImage:    String("https://www.gravatar.com/avatar/2fab21a4c4ed88e76add10650c73bae1?s=250&d=mm&r=x"),
		CoverImage:      String("a"),
		Bio:             String("b"),
		Website:         String("https://ghost.org"),
		Location:        String("The Internet"),
		Facebook:        String("ghost"),
		Twitter:         String("@ghost"),
		Accessibility:   nil,
		Status:          String("active"),
		MetaTitle:       String("a"),
		MetaDescription: String("b"),
		Tour:            Bool(true),
		LastSeen:        nil,
		CreatedAt:       Time("2019-11-26T02:39:32.000Z"),
		UpdatedAt:       Time("2019-11-26T04:30:57.000Z"),
		Roles: []*Role{
			&Role{
				ID:          String("5ddc9063c35e7700383b27e3"),
				Name:        String("Author"),
				Description: String("Authors"),
				CreatedAt:   Time("2019-11-26T02:39:31.000Z"),
				UpdatedAt:   Time("2019-11-26T02:39:31.000Z"),
			},
		},
		URL: String("https://docs.ghost.io/author/ghost-user/"),
	}
	u := &Post{
		Slug:               String("welcome-short"),
		ID:                 String("5ddc9141c35e7700383b2937"),
		UUID:               String("a5aa9bd8-ea31-415c-b452-3040dae1e730"),
		Title:              String("Welcome"),
		Mobiledoc:          String("{\"version\":\"0.3.1\",\"atoms\":[],\"cards\":[],\"markups\":[],\"sections\":[[1,\"p\",[[0,[],0,\"👋 Welcome, it's great to have you here.\"]]]]}"),
		HTML:               String("<p>👋 Welcome, it's great to have you here.</p>"),
		CommentID:          String("5ddc9141c35e7700383b2937"),
		FeatureImage:       String("https://static.ghost.org/v3.0.0/images/welcome-to-ghost.png"),
		Featured:           Bool(false),
		Status:             String("published"),
		Visibility:         String("public"),
		CreatedAt:          Time("2019-11-26T02:43:13.000Z"),
		UpdatedAt:          Time("2019-11-26T02:44:17.000Z"),
		PublishedAt:        Time("2019-11-26T02:44:17.000Z"),
		CustomExcerpt:      String("woo"),
		CodeinjectionHead:  String("<script/>"),
		CodeinjectionFoot:  String("<script/>"),
		CustomTemplate:     String("huh"),
		CanonicalURL:       String("https://pubbit.co"),
		Tags:               []*Tag{primaryTag},
		Authors:            []*Author{primaryAuthor},
		PrimaryAuthor:      primaryAuthor,
		PrimaryTag:         primaryTag,
		URL:                String("https://docs.ghost.io/welcome-short/"),
		Excerpt:            String("👋 Welcome, it's great to have you here."),
		ReadingTime:        Int(1),
		OgImage:            String("a"),
		OgTitle:            String("b"),
		OgDescription:      String("c"),
		TwitterImage:       String("d"),
		TwitterTitle:       String("e"),
		TwitterDescription: String("f"),
		MetaTitle:          String("g"),
		MetaDescription:    String("h"),
	}

	want := `{
		"slug": "welcome-short",
		"id": "5ddc9141c35e7700383b2937",
		"uuid": "a5aa9bd8-ea31-415c-b452-3040dae1e730",
		"title": "Welcome",
		"mobiledoc": "{\"version\":\"0.3.1\",\"atoms\":[],\"cards\":[],\"markups\":[],\"sections\":[[1,\"p\",[[0,[],0,\"👋 Welcome, it's great to have you here.\"]]]]}",
		"html": "<p>👋 Welcome, it's great to have you here.</p>",
		"comment_id": "5ddc9141c35e7700383b2937",
		"feature_image": "https://static.ghost.org/v3.0.0/images/welcome-to-ghost.png",
		"featured": false,
		"status": "published",
		"visibility": "public",
		"created_at": "2019-11-26T02:43:13.000Z",
		"updated_at": "2019-11-26T02:44:17.000Z",
		"published_at": "2019-11-26T02:44:17.000Z",
		"custom_excerpt": "woo",
		"codeinjection_head": "<script/>",
		"codeinjection_foot": "<script/>",
		"custom_template": "huh",
		"canonical_url": "https://pubbit.co",
		"tags": [
		  {
			"id": "5ddc9063c35e7700383b27e0",
			"name": "Getting Started",
			"slug": "getting-started",
			"description": "a",
			"feature_image": "b",
			"visibility": "public",
			"meta_title": "c",
			"meta_description": "d",
			"created_at": "2019-11-26T02:39:31.000Z",
			"updated_at": "2019-11-26T02:39:31.000Z",
			"url": "https://docs.ghost.io/tag/getting-started/"
		  }
		],
		"authors": [
		  {
			"id": "5951f5fca366002ebd5dbef7",
			"name": "Ghost",
			"slug": "ghost-user",
			"email": "info@ghost.org",
			"profile_image": "https://www.gravatar.com/avatar/2fab21a4c4ed88e76add10650c73bae1?s=250&d=mm&r=x",
			"cover_image": "a",
			"bio": "b",
			"website": "https://ghost.org",
			"location": "The Internet",
			"facebook": "ghost",
			"twitter": "@ghost",
			"accessibility": null,
			"status": "active",
			"meta_title": "a",
			"meta_description": "b",
			"tour": true,
			"last_seen": null,
			"created_at": "2019-11-26T02:39:32.000Z",
			"updated_at": "2019-11-26T04:30:57.000Z",
			"roles": [
			  {
				"id": "5ddc9063c35e7700383b27e3",
				"name": "Author",
				"description": "Authors",
				"created_at": "2019-11-26T02:39:31.000Z",
				"updated_at": "2019-11-26T02:39:31.000Z"
			  }
			],
			"url": "https://docs.ghost.io/author/ghost-user/"
		  }
		],
		"primary_author": {
		  "id": "5951f5fca366002ebd5dbef7",
		  "name": "Ghost",
		  "slug": "ghost-user",
		  "email": "info@ghost.org",
		  "profile_image": "https://www.gravatar.com/avatar/2fab21a4c4ed88e76add10650c73bae1?s=250&d=mm&r=x",
		  "cover_image": "a",
		  "bio": "b",
		  "website": "https://ghost.org",
		  "location": "The Internet",
		  "facebook": "ghost",
		  "twitter": "@ghost",
		  "accessibility": null,
		  "status": "active",
		  "meta_title": "a",
		  "meta_description": "b",
		  "tour": true,
		  "last_seen": null,
		  "created_at": "2019-11-26T02:39:32.000Z",
		  "updated_at": "2019-11-26T04:30:57.000Z",
		  "roles": [
			{
			  "id": "5ddc9063c35e7700383b27e3",
			  "name": "Author",
			  "description": "Authors",
			  "created_at": "2019-11-26T02:39:31.000Z",
			  "updated_at": "2019-11-26T02:39:31.000Z"
			}
		  ],
		  "url": "https://docs.ghost.io/author/ghost-user/"
		},
		"primary_tag": {
		  "id": "5ddc9063c35e7700383b27e0",
		  "name": "Getting Started",
		  "slug": "getting-started",
		  "description": "a",
		  "feature_image": "b",
		  "visibility": "public",
		  "meta_title": "c",
		  "meta_description": "d",
		  "created_at": "2019-11-26T02:39:31.000Z",
		  "updated_at": "2019-11-26T02:39:31.000Z",
		  "url": "https://docs.ghost.io/tag/getting-started/"
		},
		"url": "https://docs.ghost.io/welcome-short/",
		"excerpt": "👋 Welcome, it's great to have you here.",
		"reading_time": 1,
		"og_image": "a",
		"og_title": "b",
		"og_description": "c",
		"twitter_image": "d",
		"twitter_title": "e",
		"twitter_description": "f",
		"meta_title": "g",
		"meta_description": "h"
	}`

	testJSONMarshal(t, u, want)
}

func TestPostsService_Get(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(BaseAdminPath+"posts/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id": "1"}`)
	})

	post, err := client.Posts.Get("1")
	if err != nil {
		t.Errorf("Posts.Get returned error: %v", err)
	}

	want := &Post{ID: String("1")}
	if !reflect.DeepEqual(post, want) {
		t.Errorf("Posts.Get returned %+v, want %+v", post, want)
	}
}
