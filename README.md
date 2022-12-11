# Music-Library API

API for the Music-Library project, written in Go.

## Features

### Core

-   [ ] `/tracks`
-   [ ] `/albums` - array of trackIds (maybe an object with hash of album + album artist?)
-   [ ] `/track/audio`
-   [ ] `/track/cover/:size?`
-   [ ] Extract metadata
-   [ ] Cache metadata

### Additional

-   [ ] `/health`
-   [ ] `/health/metrics` // Prometheus metrics?
-   [ ] [socket.io](https://github.com/ambelovsky/gosf)
    -   [ ] Active user count (existing functionality)
    -   [ ] Session following
    -   [ ] Chat? - encrypted maybe?

### _Just for fun_

-   [ ] tests
-   [ ] benchmarks
-   [ ] n-gram search
-   [ ] file watcher, re-index after file change (wait a bit before re-indexing to avoid spamming)
-   [ ] [socket.io](https://github.com/ambelovsky/gosf)
    -   [ ] Chat? - encrypted maybe?
