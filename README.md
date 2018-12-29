[![GoDoc](https://godoc.org/github.com/im7mortal/kmutex?status.svg)](https://godoc.org/github.com/im7mortal/kmutex)
# kmutex
Sync primitive for golang. Allow block part of resource by unique ID.

It's not distributed lock.

Check [golang.org/x/sync/singleflight](https://godoc.org/golang.org/x/sync/singleflight) if you wanna reduce number of calls to the same resource. Use kmutex if you want only one caller could use resource at time.

[Kubernetes kmutex](https://godoc.org/github.com/kubernetes/kubernetes/pkg/util/keymutex) allow limit total number of calls to resource. Check implementation, it's very straight forward. You would prefer to copy it locally because it also contain kubernetes logger calls.  :heavy_exclamation_mark: That implementation can be cause `significant` delays if you don't know how it works underhood.

[GO PLAYGROUND EXAMPLE](https://play.golang.org/p/B-LBepY9rn)
