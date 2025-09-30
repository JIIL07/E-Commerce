package utils

import (
	"context"
	"runtime"
	"sync"
)

type WorkerPool struct {
	workers int
	jobs    chan func()
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewWorkerPool(workers int) *WorkerPool {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	wp := &WorkerPool{
		workers: workers,
		jobs:    make(chan func(), workers*2),
		ctx:     ctx,
		cancel:  cancel,
	}
	
	wp.start()
	return wp
}

func (wp *WorkerPool) start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	
	for {
		select {
		case job := <-wp.jobs:
			job()
		case <-wp.ctx.Done():
			return
		}
	}
}

func (wp *WorkerPool) Submit(job func()) {
	select {
	case wp.jobs <- job:
	case <-wp.ctx.Done():
	}
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.cancel()
	wp.wg.Wait()
}

func ParallelMap[T any, R any](items []T, fn func(T) R) []R {
	if len(items) == 0 {
		return []R{}
	}
	
	results := make([]R, len(items))
	var wg sync.WaitGroup
	
	for i, item := range items {
		wg.Add(1)
		go func(index int, value T) {
			defer wg.Done()
			results[index] = fn(value)
		}(i, item)
	}
	
	wg.Wait()
	return results
}

func ParallelFilter[T any](items []T, fn func(T) bool) []T {
	if len(items) == 0 {
		return []T{}
	}
	
	type result struct {
		index int
		value T
		keep  bool
	}
	
	results := make(chan result, len(items))
	var wg sync.WaitGroup
	
	for i, item := range items {
		wg.Add(1)
		go func(index int, value T) {
			defer wg.Done()
			results <- result{
				index: index,
				value: value,
				keep:  fn(value),
			}
		}(i, item)
	}
	
	go func() {
		wg.Wait()
		close(results)
	}()
	
	filtered := make([]T, 0, len(items))
	for res := range results {
		if res.keep {
			filtered = append(filtered, res.value)
		}
	}
	
	return filtered
}

func ParallelReduce[T any, R any](items []T, initial R, fn func(R, T) R, combine func(R, R) R) R {
	if len(items) == 0 {
		return initial
	}
	
	if len(items) == 1 {
		return fn(initial, items[0])
	}
	
	mid := len(items) / 2
	
	var left, right R
	var wg sync.WaitGroup
	
	wg.Add(2)
	
	go func() {
		defer wg.Done()
		left = ParallelReduce(items[:mid], initial, fn, combine)
	}()
	
	go func() {
		defer wg.Done()
		right = ParallelReduce(items[mid:], initial, fn, combine)
	}()
	
	wg.Wait()
	return combine(left, right)
}

func BatchProcess[T any](items []T, batchSize int, fn func([]T) error) error {
	if len(items) == 0 {
		return nil
	}
	
	if batchSize <= 0 {
		batchSize = runtime.NumCPU()
	}
	
	batches := make([][]T, 0, (len(items)+batchSize-1)/batchSize)
	
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}
	
	var wg sync.WaitGroup
	errChan := make(chan error, len(batches))
	
	for _, batch := range batches {
		wg.Add(1)
		go func(b []T) {
			defer wg.Done()
			if err := fn(b); err != nil {
				errChan <- err
			}
		}(batch)
	}
	
	wg.Wait()
	close(errChan)
	
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	
	return nil
}

type Future[T any] struct {
	result T
	err    error
	done   chan struct{}
}

func NewFuture[T any](fn func() (T, error)) *Future[T] {
	f := &Future[T]{
		done: make(chan struct{}),
	}
	
	go func() {
		defer close(f.done)
		f.result, f.err = fn()
	}()
	
	return f
}

func (f *Future[T]) Get() (T, error) {
	<-f.done
	return f.result, f.err
}

func (f *Future[T]) IsDone() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

func ParallelExecute[T any](tasks []func() (T, error)) ([]T, error) {
	if len(tasks) == 0 {
		return []T{}, nil
	}
	
	futures := make([]*Future[T], len(tasks))
	
	for i, task := range tasks {
		futures[i] = NewFuture(task)
	}
	
	results := make([]T, len(tasks))
	
	for i, future := range futures {
		result, err := future.Get()
		if err != nil {
			return nil, err
		}
		results[i] = result
	}
	
	return results, nil
}
