// Package uuid implements a simple, panic-free generator for a uuid v4 based
// on the AES family of block ciphers. The generator never returns an error.
package uuid

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"sync/atomic"
)

const (
	ng     = 4
	ngmask = ng - 1
)

var (
	access     [ng]uint32
	generators [ng]gen
)

func init() {
	if (ng & ngmask) != 0 {
		panic("uuid: you modified ng to create more gens, but it needs to be a power of 2")
	}
}

// V4 returns a UUIDv4. It never returns an error, never panics,
// and never runs out of entropy.
func V4() string {
	i := 0
	for {
		// this is a spinlock, nothing starves these generators
		// enough for anything more elaborate, as V4 completes
		// in 25-100ns on modern systems
		if atomic.CompareAndSwapUint32(&access[i], 0, 1) {
			u := string(generators[i].V4())
			atomic.StoreUint32(&access[i], 0)
			return u
		}
		i = (i + 1) & ngmask
	}
}

// Valid returns true if and only if the input string looks like a UUID
// It examines the length, number of dashes, their location, as well
// as the hexidecimal structure of the input string. For alphanumeric
// hex characters, both lower and upper case is allowed. It does not
// check for version bits or any other semantics not mentioned in this
// comment.
func Valid(s string) bool {
	if len(s) != 36 {
		return false
	}
	for _, n := range [...]int{8, 4, 4, 4} {
		if !acceptRun(s[:n]) {
			return false // ErrFormat
		}
		if s[n] != '-' {
			return false // ErrFieldCount
		}
		s = s[n+1:]
	}
	if !acceptRun(s) {
		return false // ErrFormat
	}
	return true
}

func acceptRun(s string) bool {
	for _, c := range []byte(s) {
		if c >= '0' && c <= '9' {
			continue
		}
		if c >= 'a' && c <= 'f' {
			continue
		}
		if c >= 'A' && c <= 'F' {
			continue
		}
		return false
	}
	return true
}

func init() {
	for i := range generators {
		g := newGen()
		generators[i] = *g
	}
}

type gen struct {
	c [36]byte
	cipher.BlockMode
	_ [8]byte
}

func newGen() *gen {
	g := &gen{}
	var k [16 * 2]byte
	_, err := rand.Read(k[:])
	if err != nil {
		panic("uuid: failed to read 32 bytes of entropy")
	}
	block, err := aes.NewCipher(k[:16])
	if err != nil {
		panic("uuid: failed to initialize aes128")
	}
	g.BlockMode = cipher.NewCBCEncrypter(block, k[16:32])
	return g
}

var h = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func (g *gen) V4() []byte {
	c := g.c[:]
	g.CryptBlocks(c[:16], c[:16])
	return append(c[:0],
		h[c[0]&15], h[c[0]>>4],
		h[c[1]&15], h[c[1]>>4],
		h[c[2]&15], h[c[2]>>4],
		h[c[3]&15], h[c[3]>>4],
		'-',
		h[c[4]&15], h[c[4]>>4],
		h[c[5]&15], h[c[5]>>4],
		'-',
		'4', h[c[6]>>4], // version 4
		h[c[7]&15], h[c[7]>>4],
		'-',
		h[8|c[8]&3], h[c[8]>>4], // clock seq hi begins with binary 10
		h[c[9]&15], h[c[9]>>4],
		'-',
		h[c[10]>>4], h[c[10]&15],
		h[c[11]>>4], h[c[11]&15],
		h[c[12]>>4], h[c[12]&15],
		h[c[13]>>4], h[c[13]&15],
		h[c[14]>>4], h[c[14]&15],
		h[c[15]>>4], h[c[15]&15],
	)
}
