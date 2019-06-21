package multilock

import (
	"sync"
)

type Multilock struct {
	mtx   sync.Mutex
	qw    map[string]*Waiter
	count map[string]int
}

func NewMultiLock() *Multilock {
	return &Multilock{qw: make(map[string]*Waiter), count: make(map[string]int)}
}

type Waiter struct {
	mtx    *sync.Mutex
	parent *Multilock
	str    string
}

type Waiters []*Waiter

func (q *Waiter) Wait() {
	q.mtx.Lock()
}

func (l *Multilock) Get(str string) *Waiter {
	l.mtx.Lock()         // Take a lock on the Multilock
	defer l.mtx.Unlock() // release it - nothing else in this function should block
	q, ok := l.qw[str]
	if !ok {
		// not in the map: make fresh, unlocked Waiter and put it there
		q = &Waiter{mtx: &sync.Mutex{}, parent: l, str: str}
		l.qw[str] = q
	}
	l.count[str] = l.count[str] + 1
	// return the Waiter in whatever state it's in
	return q
}

func (l *Multilock) GetAll(strs ...string) Waiters {
	waiters := make(Waiters, len(strs))
	for i, str := range strs {
		waiters[i] = l.Get(str)
	}
	return waiters
}

func (q *Waiter) Done() {
	q.parent.mtx.Lock()         // Take a lock on the Multilock
	defer q.parent.mtx.Unlock() // release it - nothing else in this function should block
	q, ok := q.parent.qw[q.str] // trade "this" for the reference from the parent map
	if !ok {
		panic("internal account is wrong, this shouldn't happen")
	}
	newCount := q.parent.count[q.str] - 1 // reduce the accounting
	// if no one is in line, we should clean up
	if newCount == 0 {
		delete(q.parent.qw, q.str)
		delete(q.parent.count, q.str)
		return
	}
	q.parent.count[q.str] = newCount // update accounting
	q.mtx.Unlock()                   // let anyone else who is waiting on this lock proceed
}

func (ws Waiters) Wait() {
	for _, w := range ws {
		w.Wait()
	}
}

func (ws Waiters) Done() {
	for _, w := range ws {
		w.Done()
	}
}
