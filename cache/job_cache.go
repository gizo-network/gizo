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
	MaxCacheLen = 128 //number of jobs held in cache
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

//IsFull returns true is cache is full
func (c JobCache) IsFull() bool {
	if c.getCache().Len() >= MaxCacheLen {
		return true
	}
	return false
}

//updates cache every minute
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

//Set adds key and value to cache
func (c JobCache) Set(key string, val []byte) error {
	if c.getCache().Len() >= MaxCacheLen {
		return ErrCacheFull
	}
	c.getCache().Set(key, val)
	return nil
}

//Get returns job from cache
func (c JobCache) Get(key string) (*job.Job, error) {
	jBytes, err := c.getCache().Get(key)
	if err != nil {
		return nil, err
	}
	j, err := job.DeserializeJob(jBytes)
	if err != nil {
		return nil, err
	}
	return j, nil
}

//fills up the cache with jobs with most execs in the last 15 blocks
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
		if len(sorted) > MaxCacheLen {
			for i := 0; i <= MaxCacheLen; i++ {
				c.Set(sorted[i].GetID(), sorted[i].Serialize())
			}
		} else {
			for _, job := range sorted {
				c.Set(job.GetID(), job.Serialize())
			}
		}
	} else {
		glg.Warn("Job Cache: Unable to refill cache - No blocks")
	}
}

func NewJobCache(bc *core.BlockChain) *JobCache {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(time.Minute))
	jc := JobCache{c, bc}
	jc.fill()
	// jc.watch()
	return &jc
}

// merge returns an array of job in order of number of execs in the job from max to min
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

//used to quickly sort the jobs from max execs to min execs
func mergeSort(slice []job.Job) []job.Job {
	if len(slice) < 2 {
		return slice
	}
	mid := (len(slice)) / 2
	return merge(mergeSort(slice[:mid]), mergeSort(slice[mid:]))
}
