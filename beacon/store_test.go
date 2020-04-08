package beacon

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStoreBolt(t *testing.T) {
	tmp := path.Join(os.TempDir(), "drandtest")
	require.NoError(t, os.MkdirAll(tmp, 0755))
	path := tmp
	defer os.RemoveAll(tmp)

	var sig1 = []byte{0x01, 0x02, 0x03}
	var sig2 = []byte{0x02, 0x03, 0x04}

	store, err := NewBoltStore(path, nil)
	require.NoError(t, err)

	require.Equal(t, 0, store.Len())

	b1 := &Beacon{
		PreviousSig: sig1,
		Round:       145,
		Signature:   sig2,
	}

	b2 := &Beacon{
		PreviousSig: sig2,
		Round:       146,
		Signature:   sig1,
	}

	require.NoError(t, store.Put(b1))
	require.Equal(t, 1, store.Len())
	require.NoError(t, store.Put(b1))
	require.Equal(t, 1, store.Len())
	require.NoError(t, store.Put(b2))
	require.Equal(t, 2, store.Len())

	received, err := store.Last()
	require.NoError(t, err)
	require.Equal(t, b2, received)

	store.Close()
	store, err = NewBoltStore(path, nil)
	require.NoError(t, err)
	require.NoError(t, store.Put(b1))

	doneCh := make(chan bool)
	callback := func(b *Beacon) {
		require.Equal(t, b1, b)
		doneCh <- true
	}
	store = NewCallbackStore(store, callback)
	go store.Put(b1)
	select {
	case <-doneCh:
		return
	case <-time.After(50 * time.Millisecond):
		t.Fail()
	}

	store, err = NewBoltStore(path, nil)
	require.NoError(t, err)
	store.Put(b1)
	store.Put(b2)

	store.Cursor(func(c Cursor) {
		expecteds := []*Beacon{b1, b2}
		i := 0
		for b := c.First(); b != nil; b = c.Next() {
			require.True(t, expecteds[i].Equal(b))
			i++
		}
	})
}
