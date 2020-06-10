package client

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"testing"
	"time"

	"github.com/drand/drand/chain"
)

// MockClient provide a mocked client interface
type MockClient struct {
	WatchCh chan Result
	Results []MockResult
}

// Get returns a the randomness at `round` or an error.
func (m *MockClient) Get(ctx context.Context, round uint64) (Result, error) {
	if len(m.Results) == 0 {
		return nil, errors.New("no result available")
	}
	r := m.Results[0]
	m.Results = m.Results[1:]
	return &r, nil
}

// Watch returns new randomness as it becomes available.
func (m *MockClient) Watch(ctx context.Context) <-chan Result {
	if m.WatchCh != nil {
		return m.WatchCh
	}
	ch := make(chan Result, 1)
	r, _ := m.Get(ctx, 0)
	ch <- r
	close(ch)
	return ch
}

func (m *MockClient) Info(ctx context.Context) (*chain.Info, error) {
	return nil, errors.New("not supported")
}

// RoundAt will return the most recent round of randomness
func (m *MockClient) RoundAt(_ time.Time) uint64 {
	return 0
}

// ClientWithResults returns a client on which `Get` works `m-n` times.
func MockClientWithResults(n, m uint64) *MockClient {
	c := new(MockClient)
	for i := n; i < m; i++ {
		c.Results = append(c.Results, NewMockResult(i))
	}
	return c
}

func NewMockResult(round uint64) MockResult {
	sig := make([]byte, 8)
	binary.LittleEndian.PutUint64(sig, round)
	return MockResult{
		Rnd:  round,
		Sig:  sig,
		Rand: chain.RandomnessFromSignature(sig),
	}
}

type MockResult struct {
	Rnd  uint64
	Rand []byte
	Sig  []byte
}

func (r *MockResult) Randomness() []byte {
	return r.Rand
}
func (r *MockResult) Signature() []byte {
	return r.Sig
}
func (r *MockResult) Round() uint64 {
	return r.Rnd
}
func (r *MockResult) AssertValid(t *testing.T) {
	t.Helper()
	sigTarget := make([]byte, 8)
	binary.LittleEndian.PutUint64(sigTarget, r.Rnd)
	if !bytes.Equal(r.Sig, sigTarget) {
		t.Fatalf("expected sig: %x, got %x", sigTarget, r.Sig)
	}
	randTarget := chain.RandomnessFromSignature(sigTarget)
	if !bytes.Equal(r.Rand, randTarget) {
		t.Fatalf("expected rand: %x, got %x", randTarget, r.Rand)
	}
}
