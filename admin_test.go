package ghost

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAdminClient(t *testing.T) {
	c, err := NewAdminClient("https://demo.pubbit.co", &http.Client{})
	require.NoError(t, err)
	require.NotNil(t, c.client)
}

// setup sets up a test HTTP server along with a ghost.AdminClient that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *AdminClient, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)

	// client is the GitHub client being tested and is
	// configured to use test server.
	client, err := NewAdminClient(server.URL, &http.Client{})
	if err != nil {
		log.Fatal(err)
	}

	return client, mux, server.URL, server.Close
}
