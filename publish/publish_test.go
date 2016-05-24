package main

import (
	"github.com/nats-io/nats"
	"testing"
)

func TestMissingFile(t *testing.T) {
	var urls string = nats.DefaultURL
	var sub string = "test missing file"
	var fileName string = "bar" // there is no input file with this name

	err := pub(&urls, &sub, fileName)
	if err != nil {
		t.Error(err)
	}

}
