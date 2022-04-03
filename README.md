# go-buffer
Async go-buffer that throttle when can't flush it's content.

A buffer is flushed manually or automatically, when it's full or flush interval elapsed.

`make all` for `build`, `test` and `bench`

## Buffer.Data
Adapter to use your storage provider

## Buffer.FlusherFunc
Function that called to flush buffer storage
