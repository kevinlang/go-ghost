package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/pubbit-co/go-ghost"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testName      = "test user"
	testEmail     = "test@testing.com"
	testPassword  = "testing123"
	testBlogTitle = "test blog"
)

var testBaseURL string

//TestMain controls main for the tests and allows for setup and shutdown of tests
func TestMain(m *testing.M) {
	ctx := context.Background()
	ghostC, err := getContainer(ctx)
	if err != nil {
		log.Panic(err)
	}
	defer ghostC.Terminate(ctx)
	//ghostC.FollowOutput(&logPrinter{})
	//ghostC.StartLogProducer(ctx)

	testBaseURL, err = getBaseURL(ctx, ghostC)
	if err != nil {
		log.Panic(err)
	}

	err = setupTestUser()
	if err != nil {
		log.Panic(err)
	}

	res := m.Run()
	if res != 0 {
		err = printLogs(ctx, ghostC)
		if err != nil {
			log.Panic(err)
		}
	}
	os.Exit(res)
}

func getContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "ghost",
		ExposedPorts: []string{"2368/tcp"},
		WaitingFor:   wait.ForHTTP("/").WithPort("2368"),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func getBaseURL(ctx context.Context, ghostC testcontainers.Container) (string, error) {
	host, err := ghostC.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := ghostC.MappedPort(ctx, "2368")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%s", host, port.Port()), nil
}

type logPrinter struct{}

func (*logPrinter) Accept(l testcontainers.Log) {
	log.Print(string(l.Content))
}

func printLogs(ctx context.Context, ghostC testcontainers.Container) error {
	logs, err := ghostC.Logs(ctx)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(logs)
	if err != nil {
		return err
	}

	log.Print(string(b))
	return nil
}

func setupTestUser() error {
	// we do not need any authentication for setup
	client, err := ghost.NewAdminClient(testBaseURL, &http.Client{})
	if err != nil {
		return err
	}

	details := &ghost.SetupDetails{
		Name:      testName,
		Email:     testEmail,
		Password:  testPassword,
		BlogTitle: testBlogTitle,
	}
	return client.Authentication.Setup(details)
}
