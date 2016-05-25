package main

import (
	"errors"
	"flag"
	"github.com/nats-io/nats"
	"io/ioutil"
	"log"
)

// see: https://gobyexample.com/errors
type CodedError struct {
	msg  string
	Code int // an integer code for this error
}

// CodedError implements the error interface with this
func (e *CodedError) Error() string { return e.msg }

func usage() {
	log.Fatalf("Usage: [-s server (%s)] [-sub subject] <file>\n", nats.DefaultURL)
}

// This function reads from a text file and returns the raw data of JSON
func getRaw(fpath string) ([]byte, error) {
	raw, err := ioutil.ReadFile(fpath)
	if err != nil {
		//ce := CodedError{msg: err.Error(), Code: 5}
		// return a pointer to a struct that implements the error interface:
		return nil, &CodedError{msg: err.Error(), Code: 5}
	}
	return raw, nil
}

func main() {

	// Defining command-line flags and default values
	var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
	var sub = flag.String("sub", "CVDMC", "subject to publish/subscribe on")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	var fileName string = args[0]

	_ = pub(urls, sub, fileName)
}

func pub(urls *string, sub *string, fileName string) error {

	// Connects to nats server
	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
		return errors.New("a problem with connecting to NATS")
	}
	defer nc.Close()

	msg, err := getRaw(fileName)
	if err != nil {
		return err
	}

	// Publish the data
	nc.Publish(*sub, msg)
	nc.Flush()

	log.Printf("Published [%s] : '%s'\n", *sub, msg)

	return nil
}
