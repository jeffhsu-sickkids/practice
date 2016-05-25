package main

import (
    "testing"
    "reflect"
)

func TestConvertMap(t *testing.T) {
    b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
    f := map[string]interface{}{
        "Name": "Wednesday",
        "Age":  6,
        "Parents": []interface{}{
            "Gomez",
            "Morticia",
        },
    }
    if result := convertToMap(b); reflect.DeepEqual(result, f){
        t.Error("expected", f, "got", result)
    }
}
