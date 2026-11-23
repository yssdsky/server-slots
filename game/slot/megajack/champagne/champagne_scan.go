package champagne

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

var (
	EVbot float64 // expectation of 1 bottle
	EVmjc float64 // Bottle game calculated expectation
)

func ExpBottle() {
	// avr 1 bottle gain
	EVbot = 0
	for _, v := range Bottles {
		EVbot += float64(v)
	}
	EVbot /= float64(len(Bottles))

	// expectation
	var E float64
	var n = 0
	for i := 0; i < len(Bottles); i++ {
		for j := i + 1; j < len(Bottles); j++ {
			if Bottles[i] == Bottles[j] {
				E += float64(Bottles[i]) * 4
			} else {
				E += float64(Bottles[i] + Bottles[j])
			}
			n++
		}
	}
	E /= float64(n)
	EVmjc = E
}

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	g.FSR = 15 // set free spins mode
	var s = slot.NewStatSnap(1, sn, 5, mjap, mjj)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRfg, q, _ = slot.Efs(s)
		var sq = 1 / (1 - q)
		var HRmjc = float64(s.A[0].Cb[mjc-1]) / N
		var qmjc = HRmjc / float64(g.Sel)
		var rtpmjc = EVmjc * qmjc
		var rtp = sq * (µ + rtpmjc)
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/HRfg, q, sq)
		fmt.Fprintf(w, "champagne bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRmjc, rtpmjc*100)
		if jn := s.A[0].Cj[mjj-1]; jn > 0 {
			var HRj = float64(jn) / N
			fmt.Fprintf(w, "jackpot: count %d, hit rate 1/%.8g\n", jn, 1/HRj)
		}
		fmt.Fprintf(w, "RTP = sq*(rtp(sym)+rtp(mjc)) = %.5g*(%.5g+%.5g) = %.8g%%\n", sq, µ*100, rtpmjc*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}

func CalcStatReg(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n*bonus games calculations*\n")
	ExpBottle()
	fmt.Printf("len = %d, avr bottle gain = %.5g, EV = %g\n", len(Bottles), EVbot, EVmjc)
	fmt.Printf("*free games calculations*\n")
	var rtpfs, _ = CalcStatBon(ctx, sp)
	if ctx.Err() != nil {
		return 0, 0
	}
	fmt.Printf("\n*regular games calculations*\n")
	var g = NewGame(sp.Sel)
	g.FSR = 0 // no free spins
	var s = slot.NewStatSnap(1, sn, 5, mjap, mjj)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRfg, q, _ = slot.Efs(s)
		var sq = 1 / (1 - q)
		var HRmjc = float64(s.A[0].Cb[mjc-1]) / N
		var qmjc = HRmjc / float64(g.Sel)
		var rtpmjc = EVmjc * qmjc
		var rtp = µ + rtpmjc + q*rtpfs
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/HRfg, q, sq)
		fmt.Fprintf(w, "champagne bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRmjc, rtpmjc*100)
		if jn := s.A[0].Cj[mjj-1]; jn > 0 {
			var HRj = float64(jn) / N
			fmt.Fprintf(w, "jackpot: count %d, hit rate 1/%.8g\n", jn, 1/HRj)
		}
		fmt.Fprintf(w, "RTP = rtp(sym) + rtp(mjc) + q*rtp(fg) = %.5g + %.5g + %.5g*%.5g = %.8g%%\n", µ*100, rtpmjc*100, q, rtpfs*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
