package types

var _ Locker = FackeLocker{}

type FackeLocker struct{}

func (FackeLocker) Lock()    {}
func (FackeLocker) Unlock()  {}
func (FackeLocker) RLock()   {}
func (FackeLocker) RUnlock() {}
