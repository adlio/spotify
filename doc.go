// Package spotify is a client for the Spotify Web API.
//
// The client does not handle OAuth itself. Pass New an *http.Client that
// attaches a bearer token to each request, such as one built with
// golang.org/x/oauth2. Spotify supports the Authorization Code (with or
// without PKCE) and Client Credentials flows. The implicit grant flow was
// retired on November 27, 2025.
//
// The types in this package follow the API as of the February 2026
// Development Mode changes: fields Spotify removed (popularity,
// available_markets, label, followers, linked_from) are not modeled.
package spotify
