// test package offers some common functionalities that are used throughout many
// different tests in drand.
package test

import (
	n "net"
	"strconv"

	"github.com/dedis/drand/key"
	"github.com/dedis/drand/net"
)

type testPeer struct {
	a string
	b bool
}

func (t *testPeer) Address() string {
	return t.a
}

func (t *testPeer) IsTLS() bool {
	return t.b
}

func NewPeer(addr string) net.Peer {
	return &testPeer{a: addr, b: false}
}

func NewTLSPeer(addr string) net.Peer {
	return &testPeer{a: addr, b: true}
}

// Addresses returns a list of TCP localhost addresses starting from the given
// port= start.
func Addresses(n int) []string {
	addrs := make([]string, n, n)
	for i := 0; i < n; i++ {
		addrs[i] = "127.0.0.1:" + strconv.Itoa(FreePort())
	}
	return addrs
}

// Ports returns a list of ports starting from the given
// port= start.
func Ports(n int) []string {
	ports := make([]string, n, n)
	for i := 0; i < n; i++ {
		ports[i] = strconv.Itoa(FreePort())
	}
	return ports
}

// GetFreePort returns an free TCP port.
// Taken from https://github.com/phayes/freeport/blob/master/freeport.go
func FreePort() int {
	addr, err := n.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := n.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*n.TCPAddr).Port
}

// GenerateIDs returns n keys with random port localhost addresses
func GenerateIDs(n int) []*key.Pair {
	keys := make([]*key.Pair, n)
	addrs := Addresses(n)
	for i := range addrs {
		priv := key.NewKeyPair(addrs[i])
		keys[i] = priv
	}
	return keys
}

// Generates n insecure identities
func BatchIdentities(n int) ([]*key.Pair, *key.Group) {
	privs := GenerateIDs(n)
	group := key.NewGroup(ListFromPrivates(privs), key.DefaultThreshold(n))
	return privs, group
}

// Generates n secure (TLS) identities
func BatchTLSIdentities(n int) ([]*key.Pair, *key.Group) {
	pairs, group := BatchIdentities(n)
	for i := 0; i < n; i++ {
		pairs[i].Public.TLS = true
	}
	return pairs, group
}

// ListFromPrivates returns a list of Identity from a list of Pair keys.
func ListFromPrivates(keys []*key.Pair) []*key.Identity {
	n := len(keys)
	list := make([]*key.Identity, n, n)
	for i := range keys {
		list[i] = keys[i].Public
	}
	return list

}
