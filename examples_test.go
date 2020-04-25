package ghost

import (
	"context"
	"log"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
	"golang.org/x/oauth2"
)

func ExampleNewAdminClient() {
	ts, err := NewAdminTokenSource(ExampleAdminKey)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := oauth2.NewClient(context.Background(), ts)
	client, err := NewAdminClient("https://demo.pubbit.io", httpClient)
	if err != nil {
		log.Fatal(err)
	}

	client.Posts.List(nil)
}

func ExampleNewAdminClient_Session() {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{Jar: jar}
	client, err := NewAdminClient("https://demo.pubbit.io", httpClient)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Session.Create("username", "password")
	if err != nil {
		log.Fatal(err)
	}

	client.Posts.List(nil)
}
