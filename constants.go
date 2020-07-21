package main

import (
    "time"
    "github.com/google/uuid"
)

var (
    header_chunked = []byte("Transfer-Encoding: chunked\r\n")
    trailer = []byte("\x20\xd9\x8c\xd9\x8f\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x91\xd9\x91\xd9\x91\xd9\x91\xd9\x92\xd9\x92\xd9\x92\xd9\x92\xd9\x8f\xd9\x8f\xd9\x8f\xd9\x92\xd9\x8f\xd9\x8f\r\n\r\n")
    epoch = time.Unix(0, 0)
    zeroTime time.Time
    respDown = []byte("HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\nPragma: no-cache\r\nCache-Control: no-cache, no-store, must-revalidate\r\nExpires: Thu, 01 Jan 1970 00:00:01 GMT\r\n\r\n")
    respDownLen = len(respDown)
    respUp = []byte("HTTP/1.1 204 No Content\r\nPragma: no-cache\r\nCache-Control: no-cache, no-store, must-revalidate\r\nExpires: Thu, 01 Jan 1970 00:00:01 GMT\r\n\r\n")
    respUpLen = len(respUp)
    zeroUUID uuid.UUID
    verbGetReqLen = len(makeReqBuffer(zeroUUID, false))
    verbPostReqLen = len(makeReqBuffer(zeroUUID, true))
    reqURILen = 1 + 16 * 2 + 1
)
