package static

import "embed"

// Embed static files into the binary
//go:embed images/*
var Images embed.FS
