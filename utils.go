package main

import (
    "io"
    "io/ioutil"
)

func discardBytes(r io.Reader, n int64) error {
    _, err := io.CopyN(ioutil.Discard, r, n)
    return err
}
