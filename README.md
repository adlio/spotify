# spotify

A Go client for the [Spotify Web API](https://developer.spotify.com/documentation/web-api).

```sh
go get github.com/adlio/spotify
```

The package has no dependencies outside the standard library. It's built against the API as it stands after Spotify's [February 2026 Development Mode changes](https://developer.spotify.com/documentation/web-api/tutorials/february-2026-migration-guide).

## Authentication

The client doesn't do OAuth itself. Give it an `*http.Client` that adds a bearer token to each request. [`golang.org/x/oauth2`](https://pkg.go.dev/golang.org/x/oauth2) works well for this.

For server-side catalog lookups that don't involve a user, use the Client Credentials flow:

```go
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/adlio/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	ctx := context.Background()
	cfg := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		TokenURL:     "https://accounts.spotify.com/api/token",
	}
	client := spotify.New(cfg.Client(ctx))

	track, err := client.GetTrack(ctx, "11dFghVXANMlKmJXsNCbNl", spotify.Market("US"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(track.Name, "-", track.Artists[0].Name)
}
```

To act on behalf of a user, use the Authorization Code flow with PKCE (`oauth2.Config` with `oauth2.GenerateVerifier` and `oauth2.S256ChallengeOption`). Spotify retired the implicit grant and plain-HTTP redirect URIs in November 2025. For local development, use a loopback IP like `http://127.0.0.1:8080/callback`, not `localhost`.

## Usage

```go
album, err := client.GetAlbum(ctx, "4aawyAB9vmqN3uQ7FjRGTy")
artist, err := client.GetArtist(ctx, "0TnOYISbd1XYRBk9myaseg")

results, err := client.Search(ctx, `artist:radiohead`,
	[]spotify.SearchType{spotify.SearchTypeTrack, spotify.SearchTypeAlbum},
	spotify.Limit(10), spotify.Offset(10))
```

In Development Mode, search returns at most 10 results per type (5 by default), so page with `Offset`.

## Errors and rate limits

Error responses from the API come back as `*spotify.Error`:

```go
var apiErr *spotify.Error
if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusTooManyRequests {
	// Reason is "QUOTA_EXCEEDED" when a Development Mode quota is used up
	// rather than the short-term rate limit.
	time.Sleep(apiErr.RetryAfter)
}
```

The client never retries on its own.

## Status

Early. Supported so far: Get Track, Get Album, Get Artist, and Search (tracks, albums, artists). The package doesn't wrap endpoints that Spotify has removed for new apps, such as batch lookups, recommendations, audio features, and browse.

## License

MIT
