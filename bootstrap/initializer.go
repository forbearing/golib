package bootstrap

import (
	"context"

	"golang.org/x/sync/errgroup"
)

var _initializer = new(initializer)

type initializer struct {
	fns []func() error // run init function in current goroutine.
	gos []func() error // run init function in new goroutine and receive error in channel.
}

func (i *initializer) Register(fn ...func() error) {
	i.fns = append(i.fns, fn...)
}

func (i *initializer) RegisterGo(fn ...func() error) {
	i.gos = append(i.gos, fn...)
}

func (i *initializer) Init() error {
	for j := range i.fns {
		fn := i.fns[j]
		if fn == nil {
			continue
		}
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (i *initializer) Go() error {
	g, _ := errgroup.WithContext(context.Background())
	for _, fn := range i.gos {
		if fn == nil {
			continue
		}
		g.Go(fn)

	}
	return g.Wait()
}

func Register(fn ...func() error)   { _initializer.Register(fn...) }
func RegisterGo(fn ...func() error) { _initializer.RegisterGo(fn...) }
func Init() (err error)             { return _initializer.Init() }
func Go() (err error)               { return _initializer.Go() }
