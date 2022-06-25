package main

import (
	"fmt"
	"testing"
)

func TestGenerateChart(t *testing.T) {
	images, err := LoadNowcastImages(TestDir)
	if err != nil {
		t.Fatal(err)
	}
	rain, err := images.SeekRainInfo(positionConfig)
	if err != nil {
		t.Fatal(err)
	}
	charts, err := GenerateCharts(rain)
	if err != nil {
		t.Fatal(err)
	}
	for i := range charts {
		filename := fmt.Sprintf("%s/chart%02d.png", TestDir, i)
		err = DumpImage(charts[i], filename)
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(charts) != len(images.Nowcast) {
		t.Fatalf("len(charts) != len(images)")
	}
}

func TestAnnotate(t *testing.T) {
	images, err := LoadNowcastImages(TestDir)
	if err != nil {
		t.Fatal(err)
	}
	rain, err := images.SeekRainInfo(positionConfig)
	if err != nil {
		t.Fatal(err)
	}
	charts, err := GenerateCharts(rain)
	if err != nil {
		t.Fatal(err)
	}
	img, err := images.Nowcast[0].Annotate(positionConfig, charts[0], images.Map.Image)
	if err != nil {
		t.Fatal(err)
	}

	filename := fmt.Sprintf("%s/annotated%02d.png", TestDir, 0)
	err = DumpImage(img.Image, filename)
	if err != nil {
		t.Fatal(err)
	}
}

func TestImagesAnnotate(t *testing.T) {
	images, err := LoadNowcastImages(TestDir)
	if err != nil {
		t.Fatal(err)
	}
	rain, err := images.SeekRainInfo(positionConfig)
	if err != nil {
		t.Fatal(err)
	}
	charts, err := GenerateCharts(rain)
	if err != nil {
		t.Fatal(err)
	}
	_, err = images.Annotate(positionConfig, charts)
	if err != nil {
		t.Fatal(err)
	}
}
