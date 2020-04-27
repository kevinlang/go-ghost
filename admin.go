package ghost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	BaseAdminPath = "/ghost/api/v3/admin/"
)

// An AdminClient manages communication with the Ghost Admin API
type AdminClient struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string

	Authentication *AdminAuthenticationService
	Database       *AdminDatabaseService
	Posts          *AdminPostsService
	Redirects      *AdminRedirectsService
	Session        *AdminSessionService

	// Reuse a single struct instead of allocating one for each service on the heap.
	common adminService
}

type adminService struct {
	client *AdminClient
}

// NewAdminClient returns a new client for interacting with Ghost Admin endpoints.
// baseURL should be the base admin url of the intance, in most cases taking the form
// of e.g., https://blah.pubbit.io with no trailing slash. It may additionally
// contain the subpath, but that too must omit the trailing slash.
// httpClient should handle authentication itself
func NewAdminClient(baseURL string, httpClient *http.Client) (*AdminClient, error) {
	burl, err := parseBaseURL(baseURL)
	if err != nil {
		return nil, err
	}

	// we do not currently allow specifying the version
	burl.Path += BaseAdminPath

	c := &AdminClient{client: httpClient, BaseURL: burl, UserAgent: "go-ghost"}
	c.common.client = c
	c.Authentication = (*AdminAuthenticationService)(&c.common)
	c.Database = (*AdminDatabaseService)(&c.common)
	c.Posts = (*AdminPostsService)(&c.common)
	c.Redirects = (*AdminRedirectsService)(&c.common)
	c.Session = (*AdminSessionService)(&c.common)
	return c, nil
}

func parseBaseURL(baseURL string) (*url.URL, error) {
	burl, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s as a url", baseURL)
	}

	if burl.Path != "" {
		return nil, fmt.Errorf("base url must omit the trailing slash")
	}

	return burl, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *AdminClient) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func createFormFile(w *multipart.Writer, fieldname, filename, contenttype string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contenttype)
	return w.CreatePart(h)
}

// WriteFilePart is the "callback" responsible for writing out the file of the multipart request.
type WriteFilePart func(mpw *multipart.Writer) error

// NewUploadRequest does an upload request by doing a POST against the provided path.
// It calls out to writePart to write out the principal file part of the payload,
// then populates additional multipart params provided in params.
func (c *AdminClient) NewUploadRequest(urlStr string, writePart WriteFilePart, params map[string]string) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	mp := multipart.NewWriter(buf)
	err = writePart(mp)
	if err != nil {
		return nil, err
	}
	for name, value := range params {
		mp.WriteField(name, value)
	}
	mp.Close()

	req, err := http.NewRequest("POST", u.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	req.Header.Set("Content-Type", mp.FormDataContentType())
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *AdminClient) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("received %v status from API", resp.StatusCode)
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
