package pool

import "sync"

type Resettable interface {
	Reset()
}

type Pool[T Resettable] struct {
	pool sync.Pool
}

func New[T Resettable](constructor func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return constructor()
			},
		},
	}
}

func (p *Pool[T]) Get() T {
	obj := p.pool.Get()
	if obj == nil {
		var zero T
		return zero
	}
	return obj.(T)
}

func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.pool.Put(obj)
}
