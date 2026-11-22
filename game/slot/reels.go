package slot

import (
	"math"
	"slices"
	"strconv"
	"strings"
)

// ReelAt returns symbol on the reel at true cyclic position.
// Incoming "pos" can be greater than reel length, or can be negative.
// This function is slow-performance on reels scanning.
func ReelAt(reel []Sym, pos int) Sym {
	var n = len(reel)
	return reel[(n+pos%n)%n]
}

type Reelx [][]Sym

func (rx Reelx) Clone() Reelx {
	var clone Reelx
	clone = slices.Clone(rx)
	for i, reel := range rx {
		clone[i] = slices.Clone(reel)
	}
	return clone
}

func (rx Reelx) SetPadding(sy Pos) Reelx {
	for i, reel := range rx {
		rx[i] = append(reel, reel[:sy]...)
	}
	return rx
}

func (rx Reelx) Clear() {
	clear(rx)
}

// Returns reel at given column, index is 1-based.
func (rx Reelx) Reel(col Pos) []Sym {
	return rx[col-1]
}

// Returns total number of reshuffles with padding carried.
func (rx Reelx) Reshuffles(sy Pos) uint64 {
	var res uint64 = 1
	for _, reel := range rx {
		res *= uint64(len(reel)) - uint64(sy)
	}
	return res
}

// Returns seed position on reels from array of positions on each reel.
func (rx Reelx) PosToSeed(pos []int, sy Pos) (seed uint64) {
	var m uint64 = 1
	for i := len(rx) - 1; i >= 0; i-- {
		seed += uint64(pos[i]) * m
		m *= uint64(len(rx[i])) - uint64(sy)
	}
	return
}

// Fills the array of positions on each reel from given seed.
func (rx Reelx) SeedToPos(seed uint64, pos []int, sy Pos) {
	var cur = seed
	var rl uint64
	for i := len(rx) - 1; i >= 0; i-- {
		rl = uint64(len(rx[i])) - uint64(sy)
		pos[i] = int(cur % rl)
		cur /= rl
	}
}

func (rx Reelx) String(sy Pos) string {
	var b strings.Builder
	b.WriteString("[")
	for i, reel := range rx {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strconv.Itoa(len(reel)))
	}
	b.WriteString("]")
	return b.String()
}

type ReelsMap[T any] map[float64]T

func (m ReelsMap[T]) Clear() {
	clear(m)
}

func (m ReelsMap[T]) FindClosest(mrtp float64) (val T, rtp float64) {
	rtp = -1000 // lets to get first reels from map in any case
	for p, v := range m {
		if math.Abs(mrtp-p) < math.Abs(mrtp-rtp) {
			val, rtp = v, p
		}
	}
	return
}
