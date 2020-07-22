package main

import (
    "io"
    "io/ioutil"
    "encoding/hex"
)

func discardBytes(r io.Reader, n int64) error {
    _, err := io.CopyN(ioutil.Discard, r, n)
    return err
}

type WrappedWire struct {
    chunked io.WriteCloser
    radix io.Writer
}

func NewWrappedWire(w io.Writer) *WrappedWire {
    chunked := NewChunkedWriter(w)
    radix := hex.NewEncoder(chunked)
    return &WrappedWire{chunked, radix}
}

func (ww *WrappedWire) Write(p []byte) (n int, err error) {
    return ww.radix.Write(p)
}

func (ww *WrappedWire) Close() error {
    return ww.chunked.Close()
}

type UnwrappedWire struct {
    chunked, radix io.Reader
}

func NewUnwrappedWire(r io.Reader) *UnwrappedWire {
    chunked := NewChunkedReader(r)
    radix := hex.NewDecoder(chunked)
    return &UnwrappedWire{chunked, radix}
}

func (uw *UnwrappedWire) Read(p []byte) (n int, err error) {
    return uw.radix.Read(p)
}
