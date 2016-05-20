package main

import (
    "net/http"
    "io/ioutil"
    "flag"
    "log"
    "github.com/nats-io/nats"
)

func main(){
    // Defining command-line flags and default values
    var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
    var sub = flag.String("sub", "CVDMC", "subject to publish/subscribe on")

    log.SetFlags(0)
    flag.Parse()

    // Connects to nats server
    nc, err := nats.Connect(*urls)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    res, err := http.Get("http://fhirtest.uhn.ca/baseDstu2/Patient/49028")
    if err != nil {
        panic(err)
    }
    defer res.Body.Close()
    msg, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err)
    }

    // Publish the data
    nc.Publish(*sub, msg)
    nc.Flush()

    log.Printf("Published [%s] : '%s'\n", *sub, msg)
}
