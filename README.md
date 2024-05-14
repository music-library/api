# Music-Library API

API for the Music-Library project, written in Go.

## Features

### Core

-   [x] `/tracks`
-   [x] `/albums`
-   [x] `/track/:id`
-   [x] `/track/:id/audio`
-   [x] `/track/:id/cover/:size?`
-   [x] Extract metadata
-   [x] Cache metadata
-   [x] Re-index library every X hours

### Additional

-   [x] One instance, handle multiple libraries
    -   [x] Middleware to handle library selection (via `X-Library` header)
    -   [x] Frontend to select library
        -   [x] Default library needs to be set so FE is never blocked on what to load
        -   [x] Available libraries need to be sent to the FE
        -   [x] UI changes to select/swap library
-   [x] `/health`
-   [ ] `/health/metrics` // Prometheus metrics?
-   [x] `/reindex/:password` - Refresh all metadata (without restarting the server)
-   [x] Websocket - create a custom implementation
    -   [x] Active user count (existing functionality)
    -   [x] Currently playing tracks (existing functionality)
    -   [ ] Session following
    -   [ ] Chat? - encrypted maybe?

### _Just for fun_

-   [ ] tests
-   [ ] benchmarks
-   [ ] n-gram search
-   [ ] playlists?
-   [ ] Get average + primary color of album cover
-   [ ] video support - also stream audio only via ffmpeg on-the-fly streams
-   [ ] link straight to playing a track (track # in url - only plays if no track is playing)
-   [ ] file watcher, re-index after file change (wait a bit before re-indexing to avoid spamming)
-   [ ] `/track/random`
-   [ ] `/track/search/:query` - return audio (same as `/track/:id/audio`) - Useful for searching for a song and playing it directly

## Development

### ENV

-   `HOST` - Host to run the server on (default: `localhost`)
-   `PORT` - Port to run the server on (default: `3001`)
-   `LOG_LEVEL` - Log severity level (default: `error`)
-   `LOG_FILE` - Log file (default: `DATA_DIR/music-api.log`)
-   `DATA_DIR` - Data directory to cache and store info (default: `./data`)
-   `MUSIC_DIR` - Music directory - where all your lovely music is :) (default: `./music`)

### Setup

Setup by running the following bootstrap commands:

```bash
$ go install -mod vendor github.com/go-task/task/v3/cmd/task
$ task bootstrap
```

Available tasks are in `Taskfile.yml` and use [go-task](https://taskfile.dev/#/installation). To list all available tasks, run:

```bash
$ task --list-all
```

Start development server by running:

```bash
$ task dev
```
