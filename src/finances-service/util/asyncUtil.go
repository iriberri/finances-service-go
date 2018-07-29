package util

// Starts a go-routine that read function from a channel `c` and execute
// them one by one, in a synchronized manner.
//
// Use case:
//   Safeguard some mutable state against multi-threading.
//
// Question to ask before using this:
//   Would a `mutex` not be a better solution?
//
func StartSyncDispatcher(c chan func()) {
    go func(c chan func()) {
        f := <-c
        f()
    }(c)
}
