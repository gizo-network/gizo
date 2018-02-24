package cache

import (
	"errors"
	"time"

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
}

func (c JobCache) getCache() *bigcache.BigCache {
	return c.cache
}

func (c JobCache) IsFull() bool {
	if c.getCache().Len() >= MaxCacheLen {
		return true
	}
	return false
}

func (c JobCache) Set(key string, val []byte) error {
	if c.getCache().Len() >= MaxCacheLen {
		return ErrCacheFull
	}
	c.getCache().Set(key, val)
	return nil
}

func (c JobCache) Get(key string) ([]byte, error) {
	return c.getCache().Get(key)
}

func (c JobCache) fill(blks []core.Block) {
	var jobs []job.Job
	for _, blk := range blks {
		for _, job := range blk.GetNodes() {
			jobs = append(jobs, job.GetJob())
		}
	}
	sorted := mergeSort(jobs)
	for i := 0; i <= 128; i++ {
		c.Set(sorted[i].GetID(), sorted[i].Serialize())
	}

}

func NewJobCache(bc core.BlockChain) JobCache {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	jc := JobCache{c}
	jc.fill(bc.GetBlocksWithinMinute())
	return jc
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
