package main

import (
	//"fmt"
	"github.com/nats-io/nats"
	"testing"
	"bytes"
)

func TestGetRaw(t *testing.T) {
    b := []byte(`{"menu":{"header":"SVG Viewer","items":[{"id":"Open"},{"id":"About","label":"About Adobe CVG Viewer"}]}}`)
    raw, err := getRaw("test1.json")
	if err != nil {
		t.Error(err.Error())
	}

    b = bytes.TrimSpace(b)
    raw = bytes.TrimSpace(raw)

    if (!bytes.Equal(b,raw)){
        t.Error("expected", b, "got", raw)
    }
}

// test that a missing file will be detected and produce the appropriate error
func TestMissingFile(t *testing.T) {
	var urls string = nats.DefaultURL
	var sub string = "test missing file"
	var fileName string = "bar" // there is no input file with this name

	err := pub(&urls, &sub, fileName)
	if err == nil { // there SHOULD be an error, so if there isn't the test fails
		t.Error("expected publisher to notice missing file")
	} else {
		if cerr, ok := err.(*CodedError); ok {
			if cerr.Code != 5 {
				t.Errorf("unexpected error: %s with code %d\n", cerr.msg, cerr.Code)
			}
		} else {
			t.Error("unexpected error: " + err.Error())
		}
	}
}

func TestEmptyFile(t *testing.T) {
	var urls string = nats.DefaultURL
	var sub string = "test missing file"
	var fileName string = "test2.json" // An empty file

	err := pub(&urls, &sub, fileName)
	if err == nil { // there SHOULD be an error, so if there isn't the test fails
		t.Error("expected publisher to notice empty file")
	} else {
		if cerr, ok := err.(*CodedError); ok {
			if cerr.Code != 2 {
				t.Errorf("unexpected error: %s with code %d\n", cerr.msg, cerr.Code)
			}
		} else {
			t.Error("unexpected error: " + err.Error())
		}
	}
}


func TestInvalidJson(t *testing.T) {
	var urls string = nats.DefaultURL
	var sub string = "test missing file"
	fileName := []string{"test3.json", "test4.pdf"} // A file with invalid JSON

	for _,file := range fileName {
		err := pub(&urls, &sub, file)
		if err == nil { // there SHOULD be an error, so if there isn't the test fails
			t.Error("expected publisher to notice empty file")
		} else {
			if cerr, ok := err.(*CodedError); ok {
				if cerr.Code != 6 {
					t.Errorf("unexpected error: %s with code %d\n", cerr.msg, cerr.Code)
				}
			} else {
				t.Error("unexpected error: " + err.Error())
			}
		}
	}
}
