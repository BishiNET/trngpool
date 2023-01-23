package truerand

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/AkshatM/caprice"
)

const (
	MAX_RAND_SIZE = 1024
)

type Rand struct {
	r   []int
	key string
	pos int32
	wg  sync.WaitGroup
}

// New initialize a TRNG pool.
// Fill your Random.org API Key here.
func New(apikey string) *Rand {
	return &Rand{
		key: apikey,
	}
}

func (r *Rand) refreshByGo() {
	if cap(r.r) != MAX_RAND_SIZE {
		r.r = make([]int, MAX_RAND_SIZE)
	}
	for i := 0; i < MAX_RAND_SIZE; i++ {
		r.r[i] = rand.Int()
	}
}

// It's recommended that refreshing the pool after initialization.
func (r *Rand) Refresh(isCond ...bool) {
	if len(isCond) > 0 {
		if isCond[0] {
			r.wg.Add(1)
			defer r.wg.Done()
		}
	}
	var err error
	rng := caprice.TrueRNG(r.key)
	r.r, err = rng.GenerateIntegers(MAX_RAND_SIZE, 1, 1e9, false)
	if err != nil {
		log.Println(err)
		r.refreshByGo()
		return
	}
}

func (r *Rand) Slices() []int {
	return r.r
}
func (r *Rand) Get() int {
	if cap(r.r) == 0 {
		return 1
	}
	if r.pos == int32(len(r.r)) {
		// only refresh for the first time
		if atomic.CompareAndSwapInt32(&r.pos, int32(len(r.r)), 0) {
			r.Refresh(true)
		} else {
			// Wait until refreshing.
			r.wg.Wait()
		}
	}
	return r.r[atomic.AddInt32(&r.pos, 1)-1]
}

func (r *Rand) GetN(n int) []int {
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = r.Get()
	}
	return ret
}
