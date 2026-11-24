package slotopoldeluxe

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
	"github.com/slotopol/server/game/slot/megajack/slotopol"
)

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("*bonus games calculations*\n")
	slotopol.ExpEldorado()
	slotopol.ExpMonopoly()
	fmt.Printf("*reels calculations*\n")
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, mjap, mjj)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRmje1 = s.Cb(mje1) / N
		var qmje1 = HRmje1 / float64(g.Sel)
		var rtpmje1 = slotopol.Emje * 1 * qmje1
		var HRmje3 = s.Cb(mje3) / N
		var qmje3 = HRmje3 / float64(g.Sel)
		var rtpmje3 = slotopol.Emje * 3 * qmje3
		var HRmje6 = s.Cb(mje6) / N
		var qmje6 = HRmje6 / float64(g.Sel)
		var rtpmje6 = slotopol.Emje * 6 * qmje6
		var HRmjm = s.Cb(mjm) / N
		var qmjm = HRmjm / float64(g.Sel)
		var rtpmjm = slotopol.Emjm * qmjm
		var rtp = µ + rtpmje1 + rtpmje3 + rtpmje6 + rtpmjm
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "spin1 bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRmje1, rtpmje1*100)
		fmt.Fprintf(w, "spin3 bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRmje3, rtpmje3*100)
		fmt.Fprintf(w, "spin6 bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRmje6, rtpmje6*100)
		fmt.Fprintf(w, "monopoly bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRmjm, rtpmjm*100)
		if jn := s.A[0].Cj[mjj-1]; jn > 0 {
			var HRj = float64(jn) / N
			fmt.Fprintf(w, "jackpot: count %d, hit rate 1/%.8g\n", jn, 1/HRj)
		}
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(mje) + %.5g(mjm) = %.8g%%\n", µ*100, (rtpmje1+rtpmje3+rtpmje6)*100, rtpmjm*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReels(ctx, sp, s, g, calc)
}
