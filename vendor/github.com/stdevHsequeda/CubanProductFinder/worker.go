package storeClient

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type W struct {
	ctx  context.Context
	task func(ctx context.Context)
}

func (w *W) GetArgs() context.Context {
	return w.ctx
}

func (w *W) Task(ctx context.Context) {
	section, ok := ctx.Value("section").(Section)
	if !ok {
		logrus.Warn("'section' context value not match with 'TuEnvioSection'")
		return
	}

	sc, ok := ctx.Value("sc").(*StoreClient)
	if !ok {
		logrus.Warn("'sc' context value not match with '*StoreClient'")
		return
	}
	list, err := sc.getProductsFromSection(section)
	if err != nil {
		logrus.Warn(err)
		return
	}

	for i := range list {
		err = sc.cache.AddProduct(list[i])
		if err != nil {
			logrus.Warn(err)
			continue
		}
	}
}

type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

func NewPool(maxGoroutines int) *Pool {
	p := Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w.Task(w.GetArgs())
			}
			p.wg.Done()
		}()
	}

	return &p
}

func (p *Pool) Run(w Worker) {
	p.work <- w
}

func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
