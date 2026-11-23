package fortuneteller

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

var Ecards float64

func ExpCards() {
	var sum float64
	for c1 := 1; c1 <= 4; c1++ {
		for c2 := 1; c2 <= 4; c2++ {
			for c3 := 1; c3 <= 4; c3++ {
				sum += CardsWin(c1, c2, c3)
			}
		}
	}
	Ecards = sum / 4 / 4 / 4
}

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	g.FSR = 15 // set free spins mode
	var s = slot.NewStatSnap(1, sn, 5, cbn, 0)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRcbn = s.Cb(cbn) / N
		var qcbn = HRcbn / float64(g.Sel)
		var rtpcbn = Ecards * qcbn
		var rtp = µ + rtpcbn
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "cards bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRcbn, rtpcbn*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(cards) = %.8g%%\n", µ*100, rtpcbn*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}

func CalcStatReg(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n*bonus games calculations*\n")
	ExpCards()
	fmt.Printf("total = %d, E = %g\n", 4*4*4, Ecards)
	fmt.Printf("\n*free games calculations*\n")
	var rtpfs, _ = CalcStatBon(ctx, sp)
	if ctx.Err() != nil {
		return 0, 0
	}
	fmt.Printf("\n*regular games calculations*\n")
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, cbn, 0)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRfg, q, _ = slot.Efs(s)
		var sq = 1 / (1 - q)
		var HRcbn = s.Cb(cbn) / N
		var qcbn = HRcbn / float64(g.Sel)
		var rtpcbn = Ecards * qcbn
		var rtp = µ + rtpcbn + q*rtpfs
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/HRfg, q, sq)
		fmt.Fprintf(w, "cards bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRcbn, rtpcbn*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(cards) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, rtpcbn*100, q, rtpfs*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
