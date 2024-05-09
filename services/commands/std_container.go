package commands

import (
	"context"
	"sync"

	"github.com/royiro10/cogo/models"
)

type StdListener func(line *models.StdLine)

type StdContainer struct {
	Name       string
	mu         sync.Mutex
	NotifyChan chan models.StdLine
	view       []models.StdLine
	listeners  []*StdListener
}

func NewStdContainer(name string) *StdContainer {
	sc := &StdContainer{
		Name:       name,
		NotifyChan: make(chan models.StdLine, 2),
		view:       make([]models.StdLine, 0),
		listeners:  make([]*StdListener, 0),
	}

	return sc
}

func (sc *StdContainer) Init(ctx context.Context) {
	sc.AddListener(makeViewWriterCallbak(&sc.view))

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case line := <-sc.NotifyChan:
				sc.notify(&line)
			}
		}

	}()
}

func (sc *StdContainer) View() []models.StdLine {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	return sc.view
}

func (sc *StdContainer) ViewTail(count int) []models.StdLine {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	return sc.view[len(sc.view)-count:]
}

func (sc *StdContainer) AddListener(listener *StdListener) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.listeners = append(sc.listeners, listener)
}

func (sc *StdContainer) RemoveListener(listener *StdListener) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for i, l := range sc.listeners {
		if l == listener {
			sc.listeners[i] = sc.listeners[len(sc.listeners)-1]
			sc.listeners = sc.listeners[:len(sc.listeners)-1]
		}
	}
}

func makeViewWriterCallbak(view *[]models.StdLine) *StdListener {
	mux := &sync.Mutex{}
	var viewCallback StdListener = func(line *models.StdLine) {
		mux.Lock()
		defer mux.Unlock()

		*view = append(*view, *line)
	}

	return &viewCallback
}

func (sc *StdContainer) notify(line *models.StdLine) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for _, listener := range sc.listeners {
		listenerFunc := *listener
		listenerFunc(line)
	}
}
