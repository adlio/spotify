package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Error is returned when the Web API responds with a 4xx or 5xx status.
// Use errors.As to inspect it:
//
//	var apiErr *spotify.Error
//	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusTooManyRequests {
//		time.Sleep(apiErr.RetryAfter)
//	}
type Error struct {
	// StatusCode is the HTTP status of the response.
	StatusCode int
	// Message is Spotify's description of the error, or the HTTP status
	// text if the body didn't include one.
	Message string
	// Reason is a machine-readable cause when Spotify provides one, such
	// as "QUOTA_EXCEEDED" on a 429 caused by a Development Mode quota.
	Reason string
	// RetryAfter is how long Spotify asked the caller to wait, parsed
	// from the Retry-After header. It is zero if the header was absent.
	RetryAfter time.Duration
}

func (e *Error) Error() string {
	msg := fmt.Sprintf("spotify: %d %s", e.StatusCode, e.Message)
	if e.Reason != "" {
		msg += " (" + e.Reason + ")"
	}
	return msg
}

// maxErrorBody caps how much of an error response is read.
const maxErrorBody = 64 << 10

func decodeError(resp *http.Response) error {
	e := &Error{StatusCode: resp.StatusCode}

	var body struct {
		Error struct {
			Message string `json:"message"`
			Reason  string `json:"reason"`
		} `json:"error"`
	}
	data, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody))
	if json.Unmarshal(data, &body) == nil {
		e.Message = body.Error.Message
		e.Reason = body.Error.Reason
	}
	if e.Message == "" {
		e.Message = http.StatusText(resp.StatusCode)
	}

	if secs, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && secs > 0 {
		e.RetryAfter = time.Duration(secs) * time.Second
	}
	return e
}
