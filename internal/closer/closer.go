package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

func New(signals ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}

	if len(signals) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, signals...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}

	return c
}

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func Wait() {
	globalCloser.Wait()
}

func ClosedAll() {
	globalCloser.CloseAll()
}

// add info about error while all closed
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)
		funcs := make([]func() error, 0, len(c.funcs))

		c.mu.Lock()
		copy(funcs, c.funcs)
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < len(funcs); i++ {
			if err := <-errs; err != nil {
				log.Printf("closer: error while closing: %s", err)
			}
		}
	})
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}
