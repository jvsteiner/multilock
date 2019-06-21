## Multilock

Multilock provides a simple data structure that can be used to provide fine-grained locking over different resources - doesn't matter what they are, they just need to be identified by a string.  This may be useful in cases where transactional integrity is not provided by the underlying storage layer, (think row level locking) and needs to be implemented in application space.

In some projects I have worked on, the need for such locking mechanisms existed, because the underlying database lacked an atomic read-write transaction, and so, this behavior needed to be enforced in application space.

### Install

```
$ go get -u github.com/jvsteiner/multilock
```

Or better yet, start using modules.

### Usage


Instantiate a new Multilock:

```
ml :=  NewMultiLock()

a := ml.Get("a") // lock for "a" checked out
a.Wait() // if some other goroutine was using "a", this will block
a.Done() // the lock is released, and if no other goroutines are waiting on it, deleted.

locks := ml.GetAll("b", "c", "d")
locks.Wait()
locks.Done()  // same as above, but all three locks are waited on - only when they are all acquired do we unblock.
```

