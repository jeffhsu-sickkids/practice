package main

import (
	"flag"
	"log"
	"runtime"
	"github.com/nats-io/nats"
    "gopkg.in/mgo.v2"
    //"gopkg.in/mgo.v2/bson"
    "encoding/json"
)

type Patient struct {
    Identifiers []Identifier    `json:"identifier"`
    Names []Name     `json:"name"`
    Telecoms []Telecom   `json:"telecom"`
}

type Identifier struct {
    System string   `json:"system"`
    Value string    `json:"value"`
}

type Name struct {
    Family []string     `json:"family"`
    Given []string      `json:"given"`
}

type Telecom struct {
    System string    `json:"system"`
    Value string    `json:"value"`
    Use string  `json:"use"`
}

func usage(){
    log.Fatalf("Usage:[-s server (%s)] [-sub subject]\n", nats.DefaultURL)
}

func convertToPatient(raw []byte) Patient {
    patient := Patient{}
    json.Unmarshal(raw, &patient)
    return patient
}

func insertDoc(p Patient) {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
        }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    c:= session.DB("resource").C("patient")
    err = c.Insert(&p)

    if err != nil {
        log.Fatal(err)
    }
}

func handleMsg(m *nats.Msg, i int) {
    log.Printf("[#%d] Received on [%s]:\n %s\n", i, m.Subject, string(m.Data))
    patient := convertToPatient(m.Data)
    insertDoc(patient)
}

func main(){
    var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
    var sub = flag.String("sub", "CVDMC", "subject to publish/subscribe on")

    log.SetFlags(0)
    flag.Usage = usage
    flag.Parse()

    nc, err := nats.Connect(*urls)
    if err != nil {
        log.Fatalf("Can't connect: %v\n", err)
    }

    i := 0
    nc.Subscribe(*sub, func(msg *nats.Msg) {
		i += 1
		handleMsg(msg, i)
	})

    log.Printf("Listening on [%s]\n", *sub)
    runtime.Goexit()

}
