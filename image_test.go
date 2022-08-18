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
		if rain.Rains[i].RainRange > 0. {
			t.Logf("has rain range: %f", rain.Rains[i].RainRange)
		}
	}
	t.Logf("Rainy: %v", rain.Rainy())
}
