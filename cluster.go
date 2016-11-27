package mdq

import (
	"fmt"
	"sync"
)

type Cluster interface {
	Query(qeury string) map[string]Result
}

type cluster struct {
	dbs      map[string]DB
	reporter Reporter
}

func NewCluster(dbs map[string]DB, reporter Reporter) Cluster {
	return cluster{dbs, reporter}
}

func (c cluster) Query(query string) map[string]Result {
	results := make(map[string]Result)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for n, d := range c.dbs {
		name := n
		db := d
		wg.Add(1)
		go func() {
			defer wg.Done()

			result, err := db.Query(query)
			if err != nil {
				c.reporter.Report(fmt.Errorf("[%s] %v", name, err))
				return
			}
			mu.Lock()
			results[name] = result
			mu.Unlock()
		}()
	}
	wg.Wait()
	return results
}
