package twomillionbc

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

var Eacbn float64

func ExpAcorn() {
	var sum float64
	for _, v := range Acorn {
		sum += float64(v)
	}
	Eacbn = sum / float64(len(Acorn))
}

var Edlbn float64

func ExpDiamondLion() {
	var sum float64
	for _, v := range DiamondLion {
		sum += float64(v)
	}
	Edlbn = sum / float64(len(DiamondLion))
}

func CalcStatBon(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	var g = NewGame(sp.Sel)
	g.FSR = 4 // set free spins mode
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
	fmt.Printf("\n*bonus games calculations*\n")
	ExpAcorn()
	fmt.Printf("len = %d, E = %g\n", len(Acorn), Eacbn)
	ExpDiamondLion()
	fmt.Printf("len = %d, E = %g\n", len(DiamondLion), Edlbn)
	fmt.Printf("*bonus reels calculations*\n")
	var rtpfs, _ = CalcStatBon(ctx, sp)
	if ctx.Err() != nil {
		return 0, 0
	}
	fmt.Printf("\n*regular reels calculations*\n")
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, dlbn, 0)

	var calc = func(w io.Writer) (float64, float64) {
		var N, S, _ = s.NSQ(g.Cost())
		var µ = S / N
		var HRfg, q, _ = slot.Efs(s)
		var sq = 1 / (1 - q)
		var r5l = float64(len(g.GetReels(sp.MRTP).Reel(5)) - 2)
		var HRac = 1 / r5l
		var qac = HRac
		var rtpac = Eacbn * qac
		var HRdl = s.Cb(dlbn) / N
		var qdl = HRdl / float64(g.Sel)
		var rtpdl = Edlbn * qdl
		var rtp = µ + rtpac + rtpdl + q*rtpfs
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", µ*100)
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/HRfg, q, sq)
		fmt.Fprintf(w, "acorn bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRac, rtpac*100)
		fmt.Fprintf(w, "diamond lion bonus: hit rate 1/%.5g, rtp = %.8g%%\n", 1/HRdl, rtpdl*100)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(acorn) + %.5g(dl) + %.5g*%.5g(fg) = %.8g%%\n", µ*100, rtpac*100, rtpdl*100, q, rtpfs*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
