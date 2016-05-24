package main

import (
    "net/http"
    "flag"
    "log"
    "github.com/nats-io/nats"
    "encoding/json"
)

func usage(){
    log.Fatalf("Usage: [-s server (%s)] [-sub subject] <http query >\n", nats.DefaultURL)
}

func main(){
    // Defining command-line flags and default values
    var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
    var sub = flag.String("sub", "CVDMC", "subject to publish/subscribe on")

    log.SetFlags(0)
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

    res, err := http.Get(args[0])
    if err != nil {
        panic(err)
    }
    defer res.Body.Close()
    decoder := json.NewDecoder(res.Body)

    var ma map[string]interface{}
    err = decoder.Decode(&ma)
    if err != nil {
        panic(err)
    }

    if(ma["resourceType"] == "Bundle"){
        // If resource type is bundle, need to read the entry values, and publish each of the array value individually
            if (ma["entry"] != nil){
                for _, entry := range ma["entry"].([]interface{}){
                    data := entry.(map[string]interface{})["resource"]
                    // For each data (which is in []interface), encode into byte slices in order to be published
                    msg, err := json.Marshal(data)
                    if err != nil {
                        log.Fatal(err)
                    }
                    nc.Publish(*sub, msg)
                    nc.Flush()
                    log.Printf("Published [%s] : '%s'\n", *sub, msg)
                }
            } else {
                log.Printf("The query has empty result.")
            }
    } else if (ma["resourceType"] == "Patient"){
        // If resource type is patient, simply just encode into byte slices and publish them
        msg, err := json.Marshal(ma)
        if err != nil {
            log.Fatal(err)
        }
        nc.Publish(*sub, msg)
        nc.Flush()
        log.Printf("Published [%s] : '%s'\n", *sub, msg)
    } else {
        log.Println("new resource type needs to be handle")
    }

}
