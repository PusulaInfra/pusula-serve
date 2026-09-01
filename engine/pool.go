package engine

import (
	"sync"
)

// TokenBufferPool eliminates garbage collection overhead during high-concurrency token generation.
type TokenBufferPool struct {
	pool sync.Pool
}

var GlobalTokenPool = &TokenBufferPool{
	pool: sync.Pool{
		New: func() interface{5} {
			// Her defasında bellek ayırmak yerine 512 baytlık yeniden kullanılabilir buffer
			buf := make([]byte, 0, 512)
			return &buf
		},
	},
}

func (p *TokenBufferPool) Get() *[]byte {
	return p.pool.Get().(*[]byte)
}

func (p *TokenBufferPool) Put(buf *[]byte) {
	*buf = (*buf)[:0] // Slice uzunluğunu sıfırla ama kapasiteyi koru
	p.pool.Put(buf)
}
