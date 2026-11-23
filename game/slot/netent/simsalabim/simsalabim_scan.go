package simsalabim

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

const EVne12 = 3 * 100

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	g.FSR = 10 // set free spins mode
	var s = slot.NewSnapFlat(sn, 5)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRfg, q, _ = slot.Efs(s)
		var sq = 1 / (1 - q)
		var rtp = sq * µ
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/HRfg, q, sq)
		fmt.Fprintf(w, "RTP = sq*rtp(sym) = %.5g*%.5g = %.8g%%\n", sq, µ*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}

func CalcStatReg(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n*bonus reels calculations*\n")
	var rtpfs, _ = CalcStatBon(ctx, sp)
	if ctx.Err() != nil {
		return 0, 0
	}
	fmt.Printf("\n*regular reels calculations*\n")
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, ne12, 0)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRfg, q, _ = slot.Efs(s)
		var sq = 1 / (1 - q)
		var HRne12 = s.Cb(ne12) / N
		var qne12 = HRne12 / float64(g.Sel)
		var rtpne12 = EVne12 * qne12
		var rtp = µ + rtpne12 + q*rtpfs
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/HRfg, q, sq)
		fmt.Fprintf(w, "ne12 bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRne12, rtpne12*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(ne12) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, rtpne12*100, q, rtpfs*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
