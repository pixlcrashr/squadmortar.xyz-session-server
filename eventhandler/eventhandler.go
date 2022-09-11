package eventhandler

type Delegate[S any, E any] func(sender S, event E)

type Event[S any, E any] interface {
	Add(delegate Delegate[S, E])
	Remove(delegate Delegate[S, E])
}

type EventHandler[S any, E any] interface {
	Event[S, E]
	Invoke(sender S, event E)
}

type eventHandler[S any, E any] struct {
	delegates []Delegate[S, E]
}

func (e *eventHandler[S, E]) Add(delegate Delegate[S, E]) {
	e.delegates = append(e.delegates, delegate)
}

func (e *eventHandler[S, E]) Remove(delegate Delegate[S, E]) {
	idx := e.findDelegateIndex(delegate)
	if idx == -1 {
		return
	}

	e.delegates[idx] = e.delegates[len(e.delegates)-1]
	e.delegates = e.delegates[:len(e.delegates)-1]
}

func (e *eventHandler[S, E]) findDelegateIndex(delegate Delegate[S, E]) int {
	idx := 0

	for _, d := range e.delegates {
		if &d == &delegate {
			return idx
		}

		idx++
	}

	return idx
}

func (e *eventHandler[S, E]) Invoke(sender S, event E) {
	tmpDelegates := e.delegates

	for _, d := range tmpDelegates {
		d(sender, event)
	}
}

func New[S any, E any]() EventHandler[S, E] {
	return &eventHandler[S, E]{
		make([]Delegate[S, E], 0),
	}
}
