package tools

import "sync"

// Based on https://medium.com/code-chasm/go-concurrency-pattern-worker-pool-a437117025b1
// WorkerPool is a contract for Worker Pool implementation
type WorkerPool interface {
	AddTask(task func())
	WaitForCompletion()
}

type workerPool struct {
	maxWorker   int
	queuedTaskC chan func()
	wg          *sync.WaitGroup
}

var _ WorkerPool = &workerPool{} // make the compiler check this struct implements the interface.

func NewWorkerPool(workerCount int) *workerPool {
	ret := &workerPool{
		maxWorker:   workerCount,
		queuedTaskC: make(chan func(), 100),
		wg:          &sync.WaitGroup{},
	}

	ret.wg.Add(ret.maxWorker)

	for i := 0; i < ret.maxWorker; i++ {
		go func() {
			for task := range ret.queuedTaskC {
				task()
			}

			ret.wg.Done()
		}()
	}

	return ret
}

func (s *workerPool) AddTask(task func()) {
	s.queuedTaskC <- task
}

func (s *workerPool) WaitForCompletion() {
	close(s.queuedTaskC)
	s.wg.Wait()
}
