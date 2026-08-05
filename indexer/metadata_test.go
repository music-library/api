package indexer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataExtractionFailures(t *testing.T) {
	malformedPath := filepath.Join(t.TempDir(), "malformed.mp3")
	if err := os.WriteFile(malformedPath, []byte("not valid audio metadata"), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		filePath string
	}{
		{
			name:     "missing file",
			filePath: filepath.Join(t.TempDir(), "missing.mp3"),
		},
		{
			name:     "malformed file",
			filePath: malformedPath,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rawMetadata, err := initRawMetadataFromGoLib(test.filePath)
			assert.Error(t, err)
			assert.Nil(t, rawMetadata)

			fallback := GetEmptyMetadata()
			metadata := getRawMetadataFromGoLib(fallback, test.filePath)
			assert.Same(t, fallback, metadata)
			assert.Equal(t, GetEmptyMetadata(), metadata)

			cover, mimeType := GetTrackCover(test.filePath)
			assert.Nil(t, cover)
			assert.Equal(t, "image/jpeg", mimeType)
		})
	}
}

func TestGetTrackMetadataReturnsFallbackForMissingFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "missing.mp3")

	var metadata *Metadata
	assert.NotPanics(t, func() {
		metadata = GetTrackMetadata(filePath)
	})

	assert.NotNil(t, metadata)
	assert.Equal(t, GetEmptyMetadata(), metadata)
}
