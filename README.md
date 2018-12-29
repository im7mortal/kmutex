[![GoDoc](https://godoc.org/github.com/im7mortal/kmutex?status.svg)](https://godoc.org/github.com/im7mortal/kmutex)
# kmutex
Sync primitive for golang. Allow block part of resource by unique ID.

Check [golang.org/x/sync/singleflight](https://godoc.org/golang.org/x/sync/singleflight) if you wanna reduce number of calls to the same resource. Use kmutex if you want only one caller could use resource at time.

[GO PLAYGROUND EXAMPLE](https://play.golang.org/p/B-LBepY9rn)
