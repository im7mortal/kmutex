[![GoDoc](https://godoc.org/github.com/im7mortal/kmutex?status.svg)](https://godoc.org/github.com/im7mortal/kmutex)
# kmutex
Synchronization primitive that allows locking individual resources by unique ID.

This is not a distributed lock.

See [golang.org/x/sync/singleflight](https://godoc.org/golang.org/x/sync/singleflight) if you want to reduce the number of calls to the same resource. Use kmutex if you want only one caller to use a resource at time.

[Kubernetes kmutex](https://godoc.org/k8s.io/utils/keymutex) hashes keys to a fixed set of locks, and is useful if you do not always need a separate lock for each resource.  Take a look at the implementation, it is very straight forward.

[GO PLAYGROUND EXAMPLE](https://play.golang.org/p/TPJPmW_upWO)
