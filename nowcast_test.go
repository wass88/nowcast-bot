package main

import (
	"testing"
	"time"
)

func TestDownloadableImage(t *testing.T) {
	date, err := time.Parse(time.RFC3339, "2019-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	images, err := DownloadableImage(date)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v", images)
}

func GenGif(t *testing.T) {
	PositionConfig := PositionConfig{
		Latitude:  35.6812,
		Longitude: 139.767125,
		MapID:     1,
	}
	composer := NewNowcastComposer(PositionConfig)
	composer.downloadGIF()
}
