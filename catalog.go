package spotify

import "context"

// GetTrack returns catalog information for a single track.
// Supports the Market option.
//
// API: GET /tracks/{id}
func (c *Client) GetTrack(ctx context.Context, id string, opts ...RequestOption) (*Track, error) {
	path, err := itemPath("tracks", id)
	if err != nil {
		return nil, err
	}
	var t Track
	if err := c.get(ctx, path, opts, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetAlbum returns catalog information for a single album, including the
// first page of its tracks. Supports the Market option.
//
// API: GET /albums/{id}
func (c *Client) GetAlbum(ctx context.Context, id string, opts ...RequestOption) (*Album, error) {
	path, err := itemPath("albums", id)
	if err != nil {
		return nil, err
	}
	var a Album
	if err := c.get(ctx, path, opts, &a); err != nil {
		return nil, err
	}
	return &a, nil
}

// GetArtist returns catalog information for a single artist.
//
// API: GET /artists/{id}
func (c *Client) GetArtist(ctx context.Context, id string) (*Artist, error) {
	path, err := itemPath("artists", id)
	if err != nil {
		return nil, err
	}
	var a Artist
	if err := c.get(ctx, path, nil, &a); err != nil {
		return nil, err
	}
	return &a, nil
}
