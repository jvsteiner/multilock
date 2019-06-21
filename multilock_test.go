package multilock

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLockBasic(t *testing.T) {
	l := NewMultiLock()
	a := l.Get("a")
	assert.Equal(t, 1, len(l.count))
	assert.Equal(t, 1, len(l.qw))
	a.Wait()
	a.Done()
	assert.Equal(t, 0, len(l.count))
	assert.Equal(t, 0, len(l.qw))
}

func TestLockGoroutine(t *testing.T) {
	l := NewMultiLock()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	start := time.Now()
	go doSomethingAfter(l, 30, 0, 1, 1, wg, t)
	go doSomethingAfter(l, 40, 1, 2, 0, wg, t)
	wg.Wait()
	took := time.Since(start)
	assert.True(t, took < time.Duration(100*time.Millisecond), took)
}

func doSomethingAfter(l *Multilock, waitMs, one, two, three int, wg *sync.WaitGroup, t *testing.T) {
	time.Sleep(time.Duration(waitMs) * time.Millisecond)
	assert.Equal(t, one, l.count["a"], waitMs) // after 30/40
	qw := l.Get("a")
	assert.Equal(t, two, l.count["a"], waitMs) //after 31/41
	time.Sleep(time.Duration(waitMs) * time.Millisecond)
	qw.Wait()
	qw.Done()
	assert.Equal(t, three, l.count["a"], waitMs) //after 61/81s
	wg.Done()
}

func TestLockAllGoroutine(t *testing.T) {
	l := NewMultiLock()
	var output string
	w := l.Get("a")
	wts := l.GetAll("s", "a", "b")
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		wts.Wait()
		output += "2"
		wts.Done()
		wg.Done()
	}()
	go func() {
		w.Wait()
		time.Sleep(10 * time.Millisecond)
		output += "1"
		w.Done()
		wg.Done()
	}()
	wg.Wait()
	assert.Equal(t, "12", output)
}

func BenchmarkLock(b *testing.B) {
	l := NewMultiLock()
	for i := 0; i < b.N; i++ {
		w := l.Get("str")
		w.Wait()
		w.Done()
	}
}
