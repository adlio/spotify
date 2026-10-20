package spotify

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

// SearchType is a kind of catalog item that Search can return.
type SearchType string

// Search types supported by this package.
const (
	SearchTypeAlbum  SearchType = "album"
	SearchTypeArtist SearchType = "artist"
	SearchTypeTrack  SearchType = "track"
)

// SearchResult holds one page of results per requested type. Pages for
// types that weren't requested are nil.
type SearchResult struct {
	Albums  *Page[SimpleAlbum] `json:"albums"`
	Artists *Page[Artist]      `json:"artists"`
	Tracks  *Page[Track]       `json:"tracks"`
}

// Search searches the catalog. The query supports Spotify's field filters,
// for example `track:"Paranoid Android" artist:radiohead`. Supports the
// Market, Limit and Offset options.
//
// In Development Mode, Spotify defaults to 5 results per type and allows at
// most 10. Use Offset to page through more.
//
// API: GET /search
func (c *Client) Search(ctx context.Context, query string, types []SearchType, opts ...RequestOption) (*SearchResult, error) {
	if query == "" {
		return nil, errors.New("spotify: search query must not be empty")
	}
	if len(types) == 0 {
		return nil, errors.New("spotify: at least one search type is required")
	}

	names := make([]string, len(types))
	for i, t := range types {
		names[i] = string(t)
	}
	params := func(v url.Values) {
		v.Set("q", query)
		v.Set("type", strings.Join(names, ","))
	}

	var r SearchResult
	if err := c.get(ctx, "search", append([]RequestOption{params}, opts...), &r); err != nil {
		return nil, err
	}
	return &r, nil
}
