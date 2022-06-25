package main

import (
	"fmt"
	"image/gif"
	"os"
	"testing"
)

func TestGeneateGif(t *testing.T) {
	composer, err := NewComposer(positionConfig)
	if err != nil {
		t.Fatal(err)
	}

	rainGif := composer.ComposeGif()

	filename := fmt.Sprintf("%s/rain.gif", TestDir)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gif.EncodeAll(f, &rainGif)
	if err != nil {
		t.Fatal(err)
	}
}
