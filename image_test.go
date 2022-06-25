package main

import "testing"

func TestSeekRainInfo(t *testing.T) {
	images, err := LoadNowcastImages(TestDir)
	if err != nil {
		t.Fatal(err)
	}
	rain, err := images.SeekRainInfo(positionConfig)
	if err != nil {
		t.Fatal(err)
	}
	for i := range images.Nowcast {
		if rain.Rains[i].RainFall < 0 {
			t.Fatalf("rain.Rains[%d].RainFall is negative", i)
		}
	}
	t.Logf("Rainy: %v", rain.Rainy())
}
