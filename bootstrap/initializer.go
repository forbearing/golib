package bootstrap

import "fmt"

var _initializer = new(initializer)

type initializer struct {
	fns []initFunc // run init function in current goroutine.
	gos []initFunc // run init function in new goroutine and receive error in channel.
}
type initFunc func() error

func (i *initializer) Register(fn ...initFunc) {
	i.fns = append(i.fns, fn...)
}

func (i *initializer) RegisterGo(fn ...initFunc) {
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
	errCh := make(chan error, 1)
	for j := range i.gos {
		fn := i.gos[j]
		if fn == nil {
			fmt.Println("is nil")
			continue
		}
		go func(fn initFunc) {
			errCh <- fn()
		}(fn)

	}
	return <-errCh
}

func Register(fn ...initFunc)   { _initializer.Register(fn...) }
func RegisterGo(fn ...initFunc) { _initializer.RegisterGo(fn...) }
func Init() (err error)         { return _initializer.Init() }
func Run() (err error)          { return _initializer.Init() }
func Go() (err error)           { return _initializer.Go() }
