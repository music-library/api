# Music-Library API

API for the Music-Library project, written in Go.

## Features

### Core

-   [ ] `/tracks`
-   [ ] `/albums` - array of trackIds (maybe an object with hash of album + album artist?)
-   [ ] `/track/:id`
-   [ ] `/track/:id/audio`
-   [ ] `/track/:id/cover/:size?`
-   [ ] Extract metadata
-   [ ] Cache metadata

### Additional

-   [ ] `/health`
-   [ ] `/health/metrics` // Prometheus metrics?
-   [ ] Get average + primary color of album cover
-   [ ] [socket.io](https://github.com/ambelovsky/gosf)
    -   [ ] Active user count (existing functionality)
    -   [ ] Session following

### _Just for fun_

-   [ ] tests
-   [ ] benchmarks
-   [ ] n-gram search
-   [ ] playlists?
-   [ ] file watcher, re-index after file change (wait a bit before re-indexing to avoid spamming)
-   [ ] `/track/random`
-   [ ] `/track/search/:query` - return audio (same as `/track/:id/audio`) - Useful for searching for a song and playing it directly
-   [ ] [socket.io](https://github.com/ambelovsky/gosf)
    -   [ ] Chat? - encrypted maybe?
