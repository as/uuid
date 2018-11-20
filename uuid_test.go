package uuid

import (
	"math"
	"testing"
	"time"
)

var (
	N = 1024 * 1024
)

func TestV4(t *testing.T) {
	saw := make(map[string]bool)
	t.Log(V4())
	for i := 0; i < N; i++ {
		saw[V4()] = true
	}
	t.Log(V4())
	if len(saw) != N {
		t.Fatalf("have %d distinct, want %d", len(saw), N)
	}
}

func TestRace(t *testing.T) {
	done := make(chan bool)
	defer close(done)
	hammer := func() {
		var u string
		for i := 0; i < 1024; i++ {
			u = V4()
			if u == "" {
				panic("empty uuid")
			}
		}
		select {
		case <-done:
			return
		default:
		}
		u = u
	}
	for x := 0; x < 24; x++ {
		go hammer()
	}
	var uuids = []string{}
	for i := 0; i < N; i++ {
		uuids = append(uuids, V4())
	}
	time.Sleep(time.Second * 2)
	for i := 0; i < N; i++ {
		uuids = append(uuids, V4())
	}
	checkDistribution(t, uuids...)
}

func TestProbabilityDistribution(t *testing.T) {
	var uuids = []string{}
	for i := 0; i < N; i++ {
		uuids = append(uuids, V4())
	}
	checkDistribution(t, uuids...)
}

func checkDistribution(t *testing.T, uuids ...string) {
	t.Helper()
	const expected = 1.0 / 16.0 // ~0.0625
	var ctr [256]int
	for i := 0; i < N; i++ {
		u := []byte(V4())
		if len(u) != 36 {
			t.Fatal("bad length")
		}
		for _, c := range u {
			ctr[c]++
		}
	}
	if ctr['-'] != N*36/9 {
		t.Fatal("bad dash count")
	}
	ssq := 0.0
	for _, hex := range h {
		p := float64(ctr[hex]) / float64(N*(36-4))
		ssq += p * p
	}
	t.Log(ssq)
	have := math.Round(ssq * 1000)
	want := math.Round(expected * 1000)
	if have != want {
		t.Fatalf("non-uniform probability distribution: have %f, want %f", ssq, expected)
	}
}

func BenchmarkV4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = V4()
	}
}
func TestV4Parallel(t *testing.T) {
	t.Run("AB", func(t *testing.T) {
		t.Run("A", TestV4)
		t.Run("B", TestV4)
	})
}

func BenchmarkV4Parallel(b *testing.B) {
	b.Run("X2", func(t *testing.B) {
		t.Run("A", BenchmarkV4)
		t.Run("B", BenchmarkV4)
	})
	b.Run("X5", func(t *testing.B) {
		t.Run("A", BenchmarkV4)
		t.Run("B", BenchmarkV4)
		t.Run("C", BenchmarkV4)
		t.Run("D", BenchmarkV4)
		t.Run("E", BenchmarkV4)
	})
}
