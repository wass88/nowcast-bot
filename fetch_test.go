package main

import "testing"

var TestDir string = "./forTest/"

var positionConfig PositionConfig = PositionConfig{
	X:            490,
	Y:            323,
	CursorHeight: 20,
	MapID:        10,
	Trim: TrimConfig{
		Trim:   true,
		X:      100,
		Y:      60,
		Width:  600,
		Height: 600,
	},
}

func TestDownloadImage(t *testing.T) {
	mapID := 10
	images, err := DownloadableImage(mapID)
	if err != nil {
		t.Fatal(err)
	}
	err = images.Fetch()
	if err != nil {
		t.Fatal(err)
	}
	for i := range images.Nowcast {
		if images.Nowcast[i].Image == nil {
			t.Fatalf("image.Nowcast[%d] is nil", i)
		}
	}
	if images.Map.Image == nil {
		t.Fatal("image.Map is nil")
	}
	dump := false
	if dump {
		images.Dump(TestDir)
	}
}

func TestLoadImage(t *testing.T) {
	images, err := LoadNowcastImages(TestDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(images.Nowcast) == 0 {
		t.Fatal("images.Nowcast is empty")
	}
	for i := range images.Nowcast {
		if images.Nowcast[i].Image == nil {
			t.Fatalf("image.Nowcast[%d] is nil", i)
		}
	}
	if images.Map.Image == nil {
		t.Fatal("image.Map is nil")
	}
}
