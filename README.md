# Music-Library API

API for the Music-Library project, written in Go.

## Features

### Core

-   [x] `/tracks`
-   [ ] `/albums` - array of trackIds (maybe an object with hash of album + album artist?)
-   [x] `/track/:id`
-   [ ] `/track/:id/audio`
-   [x] `/track/:id/cover/:size?`
-   [x] Extract metadata
-   [x] Cache metadata

### Additional

-   [ ] One instance, handle multiple libraries
    -   [ ] Middleware to handle library selection (maybe use a header?)
    -   [ ] `/libraries` - list libraries
    -   [ ] Frontend to select library
        -   [ ] Default library needs to be set so FE is never blocked on what to load
        -   [ ] Available libraries need to be sent to the FE
        -   [ ] UI changes to select/swap library
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

## Development

### ENV

-   `HOST` - Host to run the server on (default: `localhost`)
-   `PORT` - Port to run the server on (default: `3001`)
-   `LOG_LEVEL` - Log severity level (default: `error`)
-   `LOG_FILE` - Log file (default: `DATA_DIR/music-api.log`)
-   `DATA_DIR` - Data directory to cache and store info (default: `./data`)
-   `MUSIC_DIR` - Music directory - where all your lovely music is :) (default: `./music`)

### Setup

```bash
$ make bootstrap
```

```bash
$ make rundev
```
