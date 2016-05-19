package main

import (
    "io/ioutil"
    "github.com/nats-io/nats"
    "log"
    "flag"
    "runtime"
)

func usage(){
    log.Fatalf("Usage:[-s server (%s)] [-sub subject]\n", nats.DefaultURL)
}

func handleMsg(m *nats.Msg, i int) {
    log.Printf("[#%d] Received on [%s]:\n %s\n", i, m.Subject, string(m.Data))
    err := ioutil.WriteFile("output.json", m.Data, 0644)
    if err != nil {
          log.Fatal(err)
    }
}

func main(){
    // Setting up command-line flags and default values
    var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
    var sub = flag.String("sub", "r.CVDMC", "subject to publish/subscribe on")

    log.SetFlags(0)
    flag.Usage = usage
    flag.Parse()

    nc, err := nats.Connect(*urls)
    if err != nil {
        log.Fatalf("Can't connect: %v\n", err)
    }

    // Subscribe to the subject, i is a counter for number of message received
    i := 0
    nc.Subscribe(*sub, func(msg *nats.Msg) {
        i += 1
        handleMsg(msg, i)
    })

    log.Printf("Listening on [%s]\n", *sub)
    runtime.Goexit()

}
