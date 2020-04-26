package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/pubbit-co/go-ghost"
	"golang.org/x/net/publicsuffix"
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

	time.Sleep(time.Millisecond * 100)

	//u, _ := url.Parse("http://localhost:2369")
	//cookies := jar.Cookies(u)
	//fmt.Println(cookies)

	db, err := client.Database.Export()
	if err != nil {
		log.Fatal(err)
	}

	problems, err := client.Database.Import(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, problem := range problems {
		fmt.Println(problem)
	}
}
