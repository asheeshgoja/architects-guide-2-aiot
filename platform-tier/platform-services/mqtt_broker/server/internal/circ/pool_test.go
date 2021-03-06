package circ

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBytesPool(t *testing.T) {
	bpool := NewBytesPool(256)
	require.NotNil(t, bpool.pool)
}

func BenchmarkNewBytesPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		NewBytesPool(256)
	}
}

func TestNewBytesPoolGet(t *testing.T) {
	bpool := NewBytesPool(256)
	buf := bpool.Get()

	require.Equal(t, make([]byte, 256), buf)
	require.Equal(t, int64(1), bpool.InUse())
}

func BenchmarkBytesPoolGet(b *testing.B) {
	bpool := NewBytesPool(256)
	for n := 0; n < b.N; n++ {
		bpool.Get()
	}
}

func TestNewBytesPoolPut(t *testing.T) {
	bpool := NewBytesPool(256)
	buf := bpool.Get()
	require.Equal(t, int64(1), bpool.InUse())
	bpool.Put(buf)
	require.Equal(t, int64(0), bpool.InUse())
}

func BenchmarkBytesPoolPut(b *testing.B) {
	bpool := NewBytesPool(256)
	buf := bpool.Get()
	for n := 0; n < b.N; n++ {
		bpool.Put(buf)
	}
}
