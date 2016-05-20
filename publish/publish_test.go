package main

import (
    "testing"

)

func TestTruth(t *testing.T) {
    getRaw("data.json")
    if true != true {
        t.Error("everything I know is wrong")
    }
}
