package main

import (
	"fmt"
	"github.com/nats-io/nats"
	"testing"
)

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
			fmt.Printf("the CodedError had code %d\n", cerr.Code)
			if cerr.Code == 5 {
				fmt.Println("the program detected the missing file as expected")
			} else {
				t.Errorf("unexpected error: %s with code %d\n", cerr.msg, cerr.Code)
			}
		} else {
			t.Error("unexpected error: " + err.Error())
		}
	}

}
