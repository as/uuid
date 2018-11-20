// Package uuid implements a simple, panic-free generator for a uuid v4 based
// on the AES family of block ciphers. The generator never returns an error.
package uuid

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// V4 returns a UUIDv4 using the default generator. It never returns an error,
// never panics, and never runs out of entropy.
func V4() string {
	return string(<-ch)
}

var defaultGen = newGen()

type gen struct {
	cipher.BlockMode
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

var ch = func() chan []byte {
	ch := make(chan []byte, 32)
	go func() {
		for {
			ch <- defaultGen.V4()
		}
	}()
	return ch
}()

func (g *gen) V4() []byte {
	var c [36]byte
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
		h[c[6]&15], h[c[6]>>4],
		h[c[7]&15], h[c[7]>>4],
		'-',
		h[c[8]&15], h[c[8]>>4],
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
