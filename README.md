## analyze leveldb

This is a simple tool to analyze leveldb.

subcommand:

* kvsize: get all key/value pair size, save in sqlite database
* stats: print goleveldb stats

script tool:

* ./scripts/kvsize.py: analyze key/value size statistics

## Build 

requirement:

glibc env:
* go 1.19
* gcc

Alpine linux env:
* gcc
* go 1.19
* musl-dev
* sqlite-dev
* sqlite-static (optional, if static link)

```bash
$ go build
```

static link:

```bash
$ CGO_ENABLED=1 go build -ldflags '-extldflags "-static"'
```

