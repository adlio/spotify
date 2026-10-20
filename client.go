package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DefaultBaseURL is the root of the Spotify Web API.
const DefaultBaseURL = "https://api.spotify.com/v1/"

const defaultUserAgent = "github.com/adlio/spotify"

// Client calls the Spotify Web API. It is safe for concurrent use.
type Client struct {
	http      *http.Client
	baseURL   string
	userAgent string
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the API root, which is useful for testing against a
// local server.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/") + "/"
	}
}

// WithUserAgent sets the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		c.userAgent = ua
	}
}

// New returns a Client that sends requests through httpClient. The
// httpClient is responsible for authentication; see the package
// documentation. If httpClient is nil, http.DefaultClient is used; it sends
// no token, so Spotify will reject every request with 401.
func New(httpClient *http.Client, opts ...Option) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{
		http:      httpClient,
		baseURL:   DefaultBaseURL,
		userAgent: defaultUserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// RequestOption sets an optional query parameter on a request.
type RequestOption func(url.Values)

// Market restricts results to content available in the given ISO 3166-1
// alpha-2 country code. Pass "from_token" to use the country of the user
// the access token belongs to.
func Market(code string) RequestOption {
	return func(v url.Values) { v.Set("market", code) }
}

// Limit sets the maximum number of items to return. Spotify sets the
// allowed range per endpoint. For Search in Development Mode, the maximum
// is 10.
func Limit(n int) RequestOption {
	return func(v url.Values) { v.Set("limit", strconv.Itoa(n)) }
}

// Offset sets the index of the first item to return, for pagination.
func Offset(n int) RequestOption {
	return func(v url.Values) { v.Set("offset", strconv.Itoa(n)) }
}

// get issues a GET request for path (relative to the base URL) and decodes
// the JSON response into v.
func (c *Client) get(ctx context.Context, path string, opts []RequestOption, v any) error {
	q := url.Values{}
	for _, opt := range opts {
		opt(q)
	}
	u := c.baseURL + path
	if len(q) > 0 {
		u += "?" + q.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return decodeError(resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("spotify: decoding response from %s: %w", path, err)
	}
	return nil
}

var errEmptyID = errors.New("spotify: id must not be empty")

// itemPath builds "<collection>/<id>" with id safely escaped.
func itemPath(collection, id string) (string, error) {
	if id == "" {
		return "", errEmptyID
	}
	return collection + "/" + url.PathEscape(id), nil
}
