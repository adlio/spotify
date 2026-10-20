package spotify

import "time"

// ExternalURLs holds links to an object outside the Web API.
type ExternalURLs struct {
	Spotify string `json:"spotify"`
}

// ExternalIDs holds industry identifiers for a track or album.
type ExternalIDs struct {
	ISRC string `json:"isrc,omitempty"`
	EAN  string `json:"ean,omitempty"`
	UPC  string `json:"upc,omitempty"`
}

// Image is a cover or artist image. Height and Width are zero when Spotify
// doesn't know them.
type Image struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// Restrictions explains why content is unavailable. Reason is one of
// "market", "product" or "explicit".
type Restrictions struct {
	Reason string `json:"reason"`
}

// Copyright is a copyright statement on an album.
type Copyright struct {
	Text string `json:"text"`
	// Type is "C" for copyright or "P" for sound recording (performance)
	// copyright.
	Type string `json:"type"`
}

// Page is one page of a paginated list. Next and Previous are empty when
// there is no such page.
type Page[T any] struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []T    `json:"items"`
}

// SimpleArtist is the artist summary nested inside tracks and albums.
type SimpleArtist struct {
	ExternalURLs ExternalURLs `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

// Artist is a full artist object.
type Artist struct {
	SimpleArtist
	Images []Image `json:"images"`
	// Genres lists the artist's genres. Spotify has marked this field
	// deprecated, and it may be empty.
	Genres []string `json:"genres"`
}

// SimpleAlbum is the album summary nested inside tracks and search results.
type SimpleAlbum struct {
	// AlbumType is "album", "single" or "compilation".
	AlbumType    string       `json:"album_type"`
	TotalTracks  int          `json:"total_tracks"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Images       []Image      `json:"images"`
	Name         string       `json:"name"`
	// ReleaseDate is formatted as "2006", "2006-01" or "2006-01-02",
	// depending on ReleaseDatePrecision ("year", "month" or "day").
	ReleaseDate          string         `json:"release_date"`
	ReleaseDatePrecision string         `json:"release_date_precision"`
	Restrictions         *Restrictions  `json:"restrictions,omitempty"`
	Type                 string         `json:"type"`
	URI                  string         `json:"uri"`
	Artists              []SimpleArtist `json:"artists"`
}

// Album is a full album object, including its first page of tracks.
type Album struct {
	SimpleAlbum
	Tracks      Page[SimpleTrack] `json:"tracks"`
	Copyrights  []Copyright       `json:"copyrights"`
	ExternalIDs ExternalIDs       `json:"external_ids"`
}

// SimpleTrack is the track summary returned inside an album.
type SimpleTrack struct {
	Artists      []SimpleArtist `json:"artists"`
	DiscNumber   int            `json:"disc_number"`
	DurationMS   int            `json:"duration_ms"`
	Explicit     bool           `json:"explicit"`
	ExternalURLs ExternalURLs   `json:"external_urls"`
	Href         string         `json:"href"`
	ID           string         `json:"id"`
	// IsPlayable is only set when the request includes a Market.
	IsPlayable   *bool         `json:"is_playable,omitempty"`
	Restrictions *Restrictions `json:"restrictions,omitempty"`
	Name         string        `json:"name"`
	TrackNumber  int           `json:"track_number"`
	Type         string        `json:"type"`
	URI          string        `json:"uri"`
	IsLocal      bool          `json:"is_local"`
}

// Duration returns the track length.
func (t SimpleTrack) Duration() time.Duration {
	return time.Duration(t.DurationMS) * time.Millisecond
}

// Track is a full track object.
type Track struct {
	SimpleTrack
	Album       SimpleAlbum `json:"album"`
	ExternalIDs ExternalIDs `json:"external_ids"`
}
