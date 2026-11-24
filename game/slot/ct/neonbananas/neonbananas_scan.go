package neonbananas

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

const Ebb = 1 // bananas bonus expectation

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, lsb6, 0).MakeSym(scat2, 6)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRbb = s.Cb(bbid) / N
		var qbb = HRbb
		var rtpbb = Ebb * qbb
		var HRls1 = s.Cb(lsb1) / N
		var qls1 = HRls1 / float64(g.Sel)
		var rtpls1 = Els * 1 * qls1
		var HRls3 = s.Cb(lsb3) / N
		var qls3 = HRls3 / float64(g.Sel)
		var rtpls3 = Els * 3 * qls3
		var HRls6 = s.Cb(lsb6) / N
		var qls6 = HRls6 / float64(g.Sel)
		var rtpls6 = Els * 6 * qls6
		var rtp = µ + rtpbb + rtpls1 + rtpls3 + rtpls6
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "bananas bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRbb, rtpbb*100)
		fmt.Fprintf(w, "lucky spin1 bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRls1, rtpls1*100)
		fmt.Fprintf(w, "lucky spin3 bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRls3, rtpls3*100)
		fmt.Fprintf(w, "lucky spin6 bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRls6, rtpls6*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(bb) + %.5g(ls) = %.8g%%\n",
			µ*100, rtpbb*100, (rtpls1+rtpls3+rtpls6)*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReels(ctx, sp, s, g, calc)
}
