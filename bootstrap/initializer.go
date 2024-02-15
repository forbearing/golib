package bootstrap

var _initializer = new(initializer)

type initializer struct {
	fns []initFunc
}
type initFunc func() error

func (i *initializer) Register(fn ...initFunc) {
	i.fns = append(i.fns, fn...)
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

func Register(fn ...initFunc) { _initializer.Register(fn...) }
func Run() error              { return _initializer.Init() }
