package main

import (
	"flag"
	"log"
	"github.com/nats-io/nats"
    "io/ioutil"
)

func usage(){
    log.Fatalf("Usage: [-s server (%s)] [-sub subject] <file>\n", nats.DefaultURL)
}

// This function reads from a text file and returns the raw data of JSON
func getRaw(fpath string) []byte {
    raw, err := ioutil.ReadFile(fpath)
    if err != nil {
        panic(err)
    }
    return raw
}

func main(){

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

    // Connects to nats server
    nc, err := nats.Connect(*urls)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    msg := getRaw(args[0])

    // Publish the data
    nc.Publish(*sub, msg)
	nc.Flush()

    log.Printf("Published [%s] : '%s'\n", *sub, msg)


}
