package main

import (
	"flag"
	"github.com/nats-io/nats"
	"io/ioutil"
	"log"
	"os"
	"bytes"
	"encoding/json"
)

// see: https://gobyexample.com/errors
type CodedError struct {
	msg  string
	Code int // an integer code for this error
}

// CodedError implements the error interface with this
func (e *CodedError) Error() string { return e.msg }

func usage() {
	log.Fatalf("Usage: [-s server (%s)] [-subj subject] <file>\n", nats.DefaultURL)
}

func isJSON(b []byte) bool {
    var js map[string]interface{}
    return json.Unmarshal(b, &js) == nil
}

// This function reads from a text file and returns the raw data of JSON
func getRaw(fpath string) ([]byte, error) {
	raw, err := ioutil.ReadFile(fpath)
	if err != nil {
		// return a pointer to a struct that implements the error interface:
		switch {
			case os.IsPermission(err):
				return nil, &CodedError{msg: err.Error(), Code: 4}
				os.Exit(2)
			case os.IsNotExist(err):
				return nil, &CodedError{msg: err.Error(), Code: 5}
			default:
				return nil, &CodedError{msg: err.Error(), Code: 1}
		}
	} else if len(bytes.TrimSpace(raw)) == 0 {
		return nil, &CodedError{msg: "File is empty", Code: 2}
	} else if !isJSON(raw) {
		return nil, &CodedError{msg: "Invalid JSON", Code: 6}
	}
	return raw, nil
}

func main() {

	// Defining command-line flags and default values
	var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
	var subj = flag.String("subj", "CVDMC", "subject to publish/subscribe on")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	var fileName string = args[0]

	err := pub(urls, subj, fileName)
	if err != nil {
		panic(err)
	}
}

func pub(urls *string, subj *string, fileName string) error {

	// Connects to nats server
	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
		return &CodedError{msg: err.Error(), Code: 3}
	}
	defer nc.Close()

	msg, err := getRaw(fileName)
	if err != nil {
		return err
	}

	// Publish the data
	nc.Publish(*subj, msg)
	nc.Flush()

	log.Printf("Published [%s] : '%s'\n", *subj, msg)

	return nil
}
