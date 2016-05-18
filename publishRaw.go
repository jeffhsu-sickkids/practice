package main

import (
	"flag"
	"log"
	"github.com/nats-io/nats"
    "io/ioutil"
)

func usage(){
    log.Fatalf("Usage: [-s server (%s)] [-sub subject]\n", nats.DefaultURL)
}

func getRaw() []byte {
    raw, err := ioutil.ReadFile("data.json")
    if err != nil {
        panic(err)
    }
    return raw
}

func main(){
    var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
    var sub = flag.String("sub", "CVDMC", "subject to publish/subscribe on")

    log.SetFlags(0)
    flag.Usage = usage
    flag.Parse()

    nc, err := nats.Connect(*urls)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    msg := getRaw()

    nc.Publish(*sub, msg)
	nc.Flush()

    log.Printf("Published [%s] : '%s'\n", *sub, msg)


}