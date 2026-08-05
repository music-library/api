package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"gitlab.com/music-library/music-api/indexer"
)

func TestTrackAudioHandlerRanges(t *testing.T) {
	tempDir := t.TempDir()
	audioPath := filepath.Join(tempDir, "track.mp3")
	if err := os.WriteFile(audioPath, []byte("0123456789"), 0o600); err != nil {
		t.Fatal(err)
	}

	const (
		libraryID = "test-library"
		trackID   = "track-id"
	)
	track := &indexer.IndexTrack{
		Id:   trackID,
		Path: audioPath,
		Metadata: &indexer.Metadata{
			Artist: "Test Artist",
			Title:  "Test Track",
		},
		Stats: indexer.GetEmptyStat(),
	}

	previousIndex := indexer.MusicLibIndex
	t.Cleanup(func() {
		indexer.MusicLibIndex = previousIndex
	})
	indexer.MusicLibIndex = indexer.IndexMany{
		DefaultKey: libraryID,
		Indexes: map[string]*indexer.Index{
			libraryID: {
				Tracks:    []*indexer.IndexTrack{track},
				TracksKey: map[string]int{trackID: 0},
			},
		},
	}

	app := fiber.New()
	app.Get("/tracks/:id/audio", func(c fiber.Ctx) error {
		c.Locals("libId", libraryID)
		return TrackAudioHandler(c)
	})

	request := func(t *testing.T, rangeHeader string) (*http.Response, []byte) {
		t.Helper()

		req := httptest.NewRequest(http.MethodGet, "/tracks/"+trackID+"/audio", nil)
		if rangeHeader != "" {
			req.Header.Set(fiber.HeaderRange, rangeHeader)
		}

		resp, err := app.Test(req, fiber.TestConfig{Timeout: 0})
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if err := resp.Body.Close(); err != nil {
			t.Fatal(err)
		}
		return resp, body
	}

	validRanges := []struct {
		name         string
		rangeHeader  string
		expectedBody string
		contentRange string
	}{
		{name: "closed range", rangeHeader: "bytes=2-5", expectedBody: "2345", contentRange: "bytes 2-5/10"},
		{name: "open ended range", rangeHeader: "bytes=7-", expectedBody: "789", contentRange: "bytes 7-9/10"},
		{name: "suffix range", rangeHeader: "bytes=-4", expectedBody: "6789", contentRange: "bytes 6-9/10"},
		{name: "end beyond file", rangeHeader: "bytes=8-99", expectedBody: "89", contentRange: "bytes 8-9/10"},
		{name: "suffix larger than file", rangeHeader: "bytes=-99", expectedBody: "0123456789", contentRange: "bytes 0-9/10"},
		{name: "initial range", rangeHeader: "bytes=0-", expectedBody: "0123456789", contentRange: "bytes 0-9/10"},
	}

	for _, test := range validRanges {
		t.Run(test.name, func(t *testing.T) {
			resp, body := request(t, test.rangeHeader)

			assert.Equal(t, fiber.StatusPartialContent, resp.StatusCode)
			assert.Equal(t, "bytes", resp.Header.Get(fiber.HeaderAcceptRanges))
			assert.Equal(t, test.contentRange, resp.Header.Get(fiber.HeaderContentRange))
			assert.Equal(t, len(test.expectedBody), int(resp.ContentLength))
			assert.Equal(t, test.expectedBody, string(body))
		})
	}

	t.Run("no range", func(t *testing.T) {
		resp, body := request(t, "")

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "bytes", resp.Header.Get(fiber.HeaderAcceptRanges))
		assert.Empty(t, resp.Header.Get(fiber.HeaderContentRange))
		assert.Equal(t, int64(10), resp.ContentLength)
		assert.Equal(t, "0123456789", string(body))
	})

	invalidRanges := []struct {
		name        string
		rangeHeader string
	}{
		{name: "missing separator", rangeHeader: "bytes"},
		{name: "empty range", rangeHeader: "bytes="},
		{name: "wrong unit", rangeHeader: "items=0-1"},
		{name: "non-numeric bounds", rangeHeader: "bytes=abc-def"},
		{name: "overflow", rangeHeader: "bytes=9223372036854775808-"},
		{name: "multiple ranges", rangeHeader: "bytes=0-1,3-4"},
		{name: "reversed range", rangeHeader: "bytes=5-3"},
		{name: "zero suffix", rangeHeader: "bytes=-0"},
		{name: "start at end of file", rangeHeader: "bytes=10-"},
		{name: "missing bounds", rangeHeader: "bytes=-"},
		{name: "extra dash", rangeHeader: "bytes=1-2-3"},
		{name: "signed number", rangeHeader: "bytes=+1-2"},
	}

	for _, test := range invalidRanges {
		t.Run(test.name, func(t *testing.T) {
			resp, body := request(t, test.rangeHeader)

			assert.Equal(t, fiber.StatusRequestedRangeNotSatisfiable, resp.StatusCode)
			assert.Equal(t, "bytes", resp.Header.Get(fiber.HeaderAcceptRanges))
			assert.Equal(t, "bytes */10", resp.Header.Get(fiber.HeaderContentRange))
			assert.Greater(t, resp.ContentLength, int64(0))
			assert.NotContains(t, resp.Header.Get(fiber.HeaderContentLength), "-")
			assert.JSONEq(t, `{"message":"requested range not satisfiable","status":416}`, string(body))
		})
	}

	t.Run("range against empty file", func(t *testing.T) {
		emptyPath := filepath.Join(tempDir, "empty.mp3")
		if err := os.WriteFile(emptyPath, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		track.Path = emptyPath
		t.Cleanup(func() {
			track.Path = audioPath
		})

		resp, body := request(t, "bytes=0-")

		assert.Equal(t, fiber.StatusRequestedRangeNotSatisfiable, resp.StatusCode)
		assert.Equal(t, "bytes */0", resp.Header.Get(fiber.HeaderContentRange))
		assert.NotContains(t, resp.Header.Get(fiber.HeaderContentLength), "-")
		assert.JSONEq(t, `{"message":"requested range not satisfiable","status":416}`, string(body))
	})
}
