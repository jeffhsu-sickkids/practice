package main

import (
    "testing"
    "bytes"

)

func TestGetRaw(t *testing.T) {
    b := []byte(`{"menu":{"header":"SVG Viewer","items":[{"id":"Open"},{"id":"About","label":"About Adobe CVG Viewer"}]}}`)
    raw := getRaw("test1.json")

    b = bytes.TrimSpace(b)
    raw = bytes.TrimSpace(raw)

    if (!bytes.Equal(b,raw)){
        t.Error("expected", b, "got", raw)
    }

}
