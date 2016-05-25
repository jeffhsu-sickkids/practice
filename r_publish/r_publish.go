package main

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "log"
    "encoding/json"
    "github.com/nats-io/nats"
    "flag"
    "fmt"
)

func usage(){
    log.Fatalf("Usage:[-s server (%s)] [-subj subject] [-id docID]\n", nats.DefaultURL)
}

func readDoc(query bson.M) map[string]interface{} {
    fmt.Println(query)
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
        }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    c:= session.DB("resource").C("patient")

    var doc map[string]interface{}

    err = c.Find(query).One(&doc)
    if err != nil {
          log.Fatal(err)
      }

    return doc

}

func main() {
    // Defining command-line flags and default values
    var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
    var subj = flag.String("subj", "r.CVDMC", "subject to publish/subscribe on")
    var id = flag.String("id", "", "ID of target document in mongo")

    log.SetFlags(0)
    flag.Usage = usage
    flag.Parse()


    // Connects to nats server
    nc, err := nats.Connect(*urls)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Close()

    query := bson.M{}
    if(*id != ""){
        query = bson.M{"_id": bson.ObjectIdHex(*id)}
    }

    doc := readDoc(query)
    msg, err := json.Marshal(doc)
    if err != nil {
       log.Fatal(err)
    }

    // Publish the data
    nc.Publish(*subj, msg)
    nc.Flush()

    log.Printf("Published [%s] : '%s'\n", *subj, msg)


}
