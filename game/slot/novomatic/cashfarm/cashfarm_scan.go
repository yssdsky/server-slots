package cashfarm

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

// Bonus game expectation
// calculation see at generator/prov/novomatic/cashfarmbon.lua
const Ebon = 50

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame()
	var s = slot.NewStatSnap(slot.FallLimit, sn, 5, farmbn, 0)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var D = math.NaN()
		var cb = s.Cb(farmbn)
		var HRfarm = cb / N
		var qfarm = HRfarm
		var rtpbon = Ebon * qfarm
		var rtp = µ + rtpbon
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "farm bonus: %g, hit rate 1/%.5g, rtp = %.8g%%\n", cb, 1/HRfarm, rtpbon*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(farm) = %.8g%%\n", µ*100, rtpbon*100, rtp*100)
		slot.Print_all(w, sp, s, rtp, D)
		return rtp, D
	}

	return slot.ScanReels(ctx, sp, s, g, calc)
}
