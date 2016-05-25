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

// see: https://gobyexample.com/errors
type CodedError struct {
	msg  string
	Code int // an integer code for this error
}

// CodedError implements the error interface with this
func (e *CodedError) Error() string { return e.msg }

func usage(){
    log.Fatalf("Usage:[-s server (%s)] [-subj subject]\n", nats.DefaultURL)
}

// Convert raw data to map[string]interface
func convertToMap(raw []byte) map[string]interface{} {
	var ma map[string]interface{}
	err:= json.Unmarshal(raw, &ma)
	if err != nil {
		panic(err)
		}
	return ma

}

// Insert the map into mongo
func insertDoc(ma interface{}) {
    // Opening a session
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
        }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)

    c:= session.DB("resource").C("patient")
    err = c.Insert(&ma)

    if err != nil {
        log.Fatal(err)
    }
}

// A function to handle incoming JSON message
func handleMsg(m *nats.Msg, i int) {
    log.Printf("[#%d] Received on [%s]:\n %s\n", i, m.Subject, string(m.Data))
	ma := convertToMap(m.Data)
    insertDoc(ma)
}

func main(){
    // Setting up command-line flags and default values
    var urls = flag.String("s", nats.DefaultURL, "nats server URLs")
    var subj = flag.String("subj", "CVDMC", "subject to publish/subscribe on")

    log.SetFlags(0)
    flag.Usage = usage
    flag.Parse()

	err := sub(urls, subj)
	if err != nil {
		panic(err)
	}

}

func sub(urls *string, subj *string) error {
	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
		return &CodedError{msg: err.Error(), Code: 3}
	}

	// Subscribe to the subject, i is a counter for number of message received
	i := 0
	nc.Subscribe(*subj, func(msg *nats.Msg) {
		i += 1
		handleMsg(msg, i)
	})

	log.Printf("Listening on [%s]\n", *subj)
	runtime.Goexit()

	return nil
}
