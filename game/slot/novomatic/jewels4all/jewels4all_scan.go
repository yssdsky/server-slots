package jewels4all

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/slotopol/server/game/slot"
)

type Stat struct {
	*slot.StatSnap
	x, y slot.Pos
}

var _ slot.Simulator = (*Stat)(nil)

func NewStat(x, y slot.Pos) *Stat {
	return &Stat{
		StatSnap: slot.NewSnapFlat(sn, 5),
		x:        x, y: y,
	}
}

func (s *Stat) Clone() slot.Simulator {
	var clone = &Stat{
		StatSnap: s.StatSnap.Clone().(*slot.StatSnap),
		x:        s.x, y: s.y,
	}
	return clone
}

func (s *Stat) Clear() {
	s.StatSnap.Clear()
	s.x, s.y = 0, 0
}

func (s *Stat) Add(as slot.Simulator) {
	var op = as.(*Stat)
	s.StatSnap.Add(op.StatSnap)
}

func (s *Stat) Simulate(g slot.SlotGame, wins *slot.Wins, mrtp float64) {
	var jj = g.(*Game)
	var sym = jj.At(s.x, s.y)
	jj.SetSym(s.x, s.y, wild)
	s.StatSnap.Simulate(g, wins, mrtp)
	jj.SetSym(s.x, s.y, sym)
}

func CalcStatEuro(ctx context.Context, sp *slot.ScanPar, s slot.Simulator) (float64, float64) {
	var g = NewGame(sp.Sel)

	var calc = func(w io.Writer) (float64, float64) {
		return slot.Parsheet_simple(w, sp, s, g.Cost())
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}

// custom parsheet
func CalcStat(ctx context.Context, sp *slot.ScanPar) (rtp, D float64) {
	var Pw, _ = ChanceMap.FindClosest(sp.MRTP) // wild chance

	var Σµ, Σµ2, ΣD float64
	for i := range 15 {
		var x = slot.Pos(i/3 + 1)
		var y = slot.Pos(i%3 + 1)
		fmt.Printf("\n(%d/16) calculations of euro at [%d,%d]\n", i+1, x, y)
		var s = NewStat(x, y)
		var µ, D = CalcStatEuro(ctx, sp, s)
		Σµ += µ
		Σµ2 += µ * µ
		ΣD += D
	}
	var µw = Σµ / 15
	var Dw = ΣD/15 + Σµ2/15 - µw*µw
	fmt.Printf("\n(16/16) regular games calculations\n")
	var sr = slot.NewSnapFlat(sn, 5)
	var µr, Dr = CalcStatEuro(ctx, sp, sr)
	rtp = (1-Pw)*µr + Pw*µw // Pw = (rtp-µr)/(µw - µr)
	D = (1-Pw)*Dr + Pw*Dw + Pw*(1-Pw)*(µw-µr)*(µw-µr)
	if sp.IsMain() {
		fmt.Printf("\n*final calculations*\n")
		fmt.Printf("euro avr: rtp(eu) = %.8g%%\n", µw*100)
		fmt.Printf("wild hit rate: 1/%.5g\n", 1/Pw)
		fmt.Printf("RTP = (1-Pw)*%.5g(sym) + Pw*%.5g(eu) = %.8g%%\n", µr*100, µw*100, rtp*100)
	}
	var w = os.Stdout
	slot.Print_all(w, sp, sr, rtp, D)
	return
}
