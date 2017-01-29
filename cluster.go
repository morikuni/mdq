package mdq

import (
	"sync"
)

type Cluster interface {
	Query(qeury string) []Result
}

type cluster struct {
	dbs      []DB
	reporter Reporter
}

func NewCluster(dbs []DB, reporter Reporter) Cluster {
	return cluster{dbs, reporter}
}

func (c cluster) Query(query string) []Result {
	var results []Result
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, d := range c.dbs {
		db := d
		wg.Add(1)
		go func() {
			defer wg.Done()

			result, err := db.Query(query)
			if err != nil {
				c.reporter.Report(err)
				return
			}
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}()
	}
	wg.Wait()
	return results
}
