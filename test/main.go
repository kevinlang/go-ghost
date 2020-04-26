package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"

	"github.com/pubbit-co/go-ghost"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/oauth2"
)

func main() {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{
		Jar: jar,
		//Transport: &ghost.AdminSessionTransport{Origin: "https://test.com"},
	}
	client, err := ghost.NewAdminClient("http://localhost:2369", httpClient)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Session.Create("testing@testing.com", "testing123")
	if err != nil {
		log.Fatal(err)
	}

	//time.Sleep(time.Millisecond * 100)

	//u, _ := url.Parse("http://localhost:2369")
	//cookies := jar.Cookies(u)
	//fmt.Println(cookies)

	db, err := client.Database.Export()
	if err != nil {
		log.Fatal(err)
	}

	ts, err := ghost.NewAdminTokenSource("5ea1aeb17edc2650468b6554:0f1103f5af0395a73041457eb6928f9e0d143a8dcba187915342e65687e2a589")
	if err != nil {
		log.Fatal(err)
	}
	keyHttpClient := oauth2.NewClient(context.Background(), ts)
	keyClient, err := ghost.NewAdminClient("http://localhost:2369", keyHttpClient)
	problems, err := keyClient.Database.Import(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, problem := range problems {
		fmt.Println(problem)
	}

	/*
		redirects, err := client.Redirects.Download()
		if err != nil {
			log.Fatal(err)
		}

		redirects = append(redirects, &ghost.Redirect{
			From: "/hello",
			To:   "/",
		})
		err = client.Redirects.Upload(redirects)
		if err != nil {
			log.Fatal(err)
		}*/
}
