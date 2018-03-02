package cache

import (
	"errors"
	"time"

	"github.com/kpango/glg"

	"github.com/allegro/bigcache"
	"github.com/gizo-network/gizo/core"
	"github.com/gizo-network/gizo/job"
)

const (
	MaxCacheLen = 128
)

var (
	ErrCacheFull = errors.New("Cache: Cache filled up")
)

type JobCache struct {
	cache *bigcache.BigCache
	bc    *core.BlockChain
}

func (c JobCache) getCache() *bigcache.BigCache {
	return c.cache
}

func (c JobCache) getBC() *core.BlockChain {
	return c.bc
}

func (c JobCache) IsFull() bool {
	if c.getCache().Len() >= MaxCacheLen {
		return true
	}
	return false
}

func (c JobCache) watch() {
	ticker := time.NewTicker(time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				glg.Warn("Job Cache: Updating cache")
				c.fill()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c JobCache) set(key string, val []byte) error {
	if c.getCache().Len() >= MaxCacheLen {
		return ErrCacheFull
	}
	c.getCache().Set(key, val)
	return nil
}

func (c JobCache) Get(key string) ([]byte, error) {
	return c.getCache().Get(key)
}

func (c JobCache) fill() {
	var jobs []job.Job
	blks := c.getBC().GetLatest15()
	if len(blks) != 0 {
		for _, blk := range blks {
			for _, job := range blk.GetNodes() {
				jobs = append(jobs, job.GetJob())
			}
		}
		sorted := mergeSort(jobs)
		if len(sorted) > 128 {
			for i := 0; i <= 128; i++ {
				c.set(sorted[i].GetID(), sorted[i].Serialize())
			}
		} else {
			for _, job := range sorted {
				c.set(job.GetID(), job.Serialize())
			}
		}
	} else {
		glg.Warn("Job Cache: Unable to refill cache - Not blocks in the last minute")
	}
}

func NewJobCache(bc *core.BlockChain) *JobCache {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	jc := JobCache{c, bc}
	jc.fill()
	jc.watch()
	return &jc
}

// Merge returns an array of job in order of number of execs in the job from max to min
func merge(left, right []job.Job) []job.Job {
	size, i, j := len(left)+len(right), 0, 0
	slice := make([]job.Job, size, size)
	for k := 0; k < size; k++ {
		if i > len(left)-1 && j <= len(right)-1 {
			slice[k] = right[j]
			j++
		} else if j > len(right)-1 && i <= len(left)-1 {
			slice[k] = left[i]
			i++
		} else if len(left[i].GetExecs()) > len(right[j].GetExecs()) {
			slice[k] = left[i]
			i++
		} else {
			slice[k] = right[j]
			j++
		}
	}
	return slice
}

func mergeSort(slice []job.Job) []job.Job {
	if len(slice) < 2 {
		return slice
	}
	mid := (len(slice)) / 2
	return merge(mergeSort(slice[:mid]), mergeSort(slice[mid:]))
}
