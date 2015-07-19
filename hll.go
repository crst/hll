package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
)

func hash(s string) (ret uint64) {
	h := sha1.New()
	h.Write([]byte(s))
	b := h.Sum(nil)
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func p(s uint64) int8 {
	r := int8(0)
	for (s & 1) != 1 {
		s >>= 1
		r += 1
	}
	return r + 1
}

func max(a int8, b int8) int8 {
	if a < b {
		return b
	}
	return a
}

func alpha(m uint64) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	}
	return 0.7213 / (1.0 + 1.079/float64(m))
}

func hll(file string, b uint64) (float64, int) {
	var m uint64 = 1 << b
	var M = make([]int8, m)
	var bs uint64 = 0
	for i := uint64(0); i < b; i++ {
		bs |= (1 << i)
	}

	var exact = make(map[int64]bool)

	f, _ := os.Open(file)
	defer f.Close()
	z, _ := gzip.NewReader(f)
	defer z.Close()
	scanner := bufio.NewScanner(z)
	for scanner.Scan() {
		line := scanner.Text()
		x := hash(line)
		j := (x & bs)
		w := x >> b
		M[j] = max(M[j], p(w))

		number, _ := strconv.ParseInt(line, 10, 64)
		exact[number] = true
	}

	Z, V := 0.0, 0.0
	for i := 0; i < int(m); i++ {
		if M[i] == 0 {
			V += 1
		}
		Z += 1.0 / float64(int(1)<<uint(M[i]))
	}

	var E = (alpha(m) * float64(m) * float64(m)) / Z

	if E < 5.0/2.0*float64(m) && V != 0 {
		E = float64(m) * math.Log2(float64(m)/V)
	}

	return E, len(exact)
}

func main() {
	var f = os.Args[1]
	var b uint64 = 8

	var est, ext = hll(f, b)
	fmt.Printf("Estimated value: %.3f\n", est)
	fmt.Printf("Exact value: %d\n", ext)
	fmt.Printf("Approximation: %.3f\n", est/float64(ext))
}
