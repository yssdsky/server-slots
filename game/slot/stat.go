package slot

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/slotopol/server/game"
)

type Simulator interface {
	// Clone returns a full deep copy of the statistics.
	Clone() Simulator
	// Clear resets the statistics to zero.
	Clear()
	// Add updates the statistics with given simulator data with same type.
	// It helps to parallelize the simulation and then merge the results.
	Add(s Simulator)
	// Returns free games hit rate, expectation of initial freespins,
	// second moment of initial freespins calculated by dedicated counters.
	Simulate(SlotGame, *Wins, float64)

	MathSimulator
}

type MathSimulator interface {
	// Count returns number of processed grid reshuffles, including with no wins.
	Count() float64
	// NSQ returns number of spins, sum of pays, sum of squares of pays by spin,
	// normalized by spin cost.
	NSQ(cost float64) (N, S, Q float64)
	// FGC returns counter of triggering free game series, total count of free spins,
	// and sum of squared free spin series lengths.
	FGC() (Cfg, Cfs, Qfs float64)
}

type ScanPar = game.ScanPar

// Returns plain math expectation (µ = S/N) and plain dispersion (D = Q/N - µ*µ).
func EvD(s MathSimulator, cost float64) (float64, float64) {
	var N, S, Q = s.NSQ(cost)
	var µ = S / N
	var D = Q/N - µ*µ
	return µ, D
}

// Returns free games hit rate, expectation of initial freespins,
// second moment of initial freespins calculated by dedicated counters.
func Efs(s MathSimulator) (HRfg, EK, EK2 float64) {
	var N = s.Count()
	var Cfg, Cfs, Qfs = s.FGC()
	HRfg = Cfg / N
	EK = Cfs / N
	EK2 = Qfs / N
	return
}

type (
	CalcFunc = func(io.Writer) (float64, float64)
	CalcAlg  = func(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, calc CalcFunc)
)

const (
	FallLimit = 20 // maximum 20 cascades for reels ~100 symbols length
)

var (
	ErrAvalanche = errors.New("too many cascading falls")
	ErrReelCount = errors.New("unexpected number of reels")
)

func ScanReels(ctx context.Context,
	sp *ScanPar, ss Simulator, g SlotGeneric, calc CalcFunc) (float64, float64) {
	var reels = g.GetReels(sp.MRTP)
	var sx, sy = g.Dim()
	if len(reels) != int(sx) {
		panic(fmt.Errorf("%w: %d reels provided for %dx%d slot", ErrReelCount, len(reels), sx, sy))
	}
	switch sp.Method {
	case game.CMbruteforce:
		sp.Total = reels.Reshuffles(sy)
		BruteForcex(ctx, sp, ss, g, calc)
	case game.CMmontecarlo:
		MonteCarlo(ctx, sp, ss, g, calc)
	}
	return calc(io.Discard)
}
