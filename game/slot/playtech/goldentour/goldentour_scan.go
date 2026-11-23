package goldentour

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

var Egolfbn float64

func ExpGolf() {
	var sum float64
	for _, v := range Golf {
		sum += float64(v)
	}
	Egolfbn = sum / float64(len(Golf))
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n*bonus games calculations*\n")
	ExpGolf()
	fmt.Printf("len = %d, E = %g\n", len(Golf), Egolfbn)
	fmt.Printf("\n*reels calculations*\n")
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, golfbon, 0)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRgolf = s.Cb(golfbon) / N
		var qgolf = HRgolf / float64(g.Sel)
		var rtpgolf = Egolfbn * qgolf
		var rtp = µ + rtpgolf
		fmt.Fprintf(w, "golf bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRgolf, rtpgolf*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(golf) = %.8g%%\n", µ*100, rtpgolf*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
