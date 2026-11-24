package fairyqueen

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/slotopol/server/game/slot"
)

func BruteForce5x3es2(ctx context.Context, sp *slot.ScanPar, s slot.Simulator, g *Game, reels slot.Reelx) {
	var tn64 = uint64(sp.TN)
	var r3 = reels.Reel(3)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)

	var final = make(chan slot.Simulator)
	defer close(final)

	for ti := range tn64 {
		var local = s.Clone()
		var wgame = g.Clone()
		var N uint64
		go func() {
			var wins slot.Wins
			for i3 := range r3[:len(r3)-2] {
				wgame.SetCol(3, r3, i3)
				for i4 := range r4[:len(r4)-2] {
					wgame.SetCol(4, r4, i4)
					for i5 := range r5[:len(r5)-2] {
						N++
						if N&0xFFF == 0 {
							select {
							case <-ctx.Done():
								return
							default:
							}
						}
						if N%tn64 != ti {
							continue
						}
						wgame.SetCol(5, r5, i5)
						local.Simulate(wgame, &wins, sp.MRTP)
						wins.Reset()
					}
				}
			}
			final <- local
		}()
	}
	for range tn64 {
		select {
		case <-ctx.Done():
			return
		case local := <-final:
			s.Add(local)
		}
	}
}

func BruteForce5x3es3(ctx context.Context, sp *slot.ScanPar, s slot.Simulator, g *Game, reels slot.Reelx) {
	var tn64 = uint64(sp.TN)
	var r4 = reels.Reel(4)
	var r5 = reels.Reel(5)

	var final = make(chan slot.Simulator)
	defer close(final)

	for ti := range tn64 {
		var local = s.Clone()
		var wgame = g.Clone().(*Game)
		var N uint64
		go func() {
			var wins slot.Wins
			for i4 := range r4[:len(r4)-2] {
				wgame.SetCol(4, r4, i4)
				for i5 := range r5[:len(r5)-2] {
					N++
					if N&0xFFF == 0 {
						select {
						case <-ctx.Done():
							return
						default:
						}
					}
					if N%tn64 != ti {
						continue
					}
					wgame.SetCol(5, r5, i5)
					local.Simulate(wgame, &wins, sp.MRTP)
					wins.Reset()
				}
			}
			final <- local
		}()
	}
	for range tn64 {
		select {
		case <-ctx.Done():
			return
		case local := <-final:
			s.Add(local)
		}
	}
}

func CalcStatBon(ctx context.Context, sp *slot.ScanPar, es slot.Sym) (float64, float64, float64) {
	var g = NewGame(sp.Sel)
	g.FSR = 10 // set free spins mode
	g.ES = es
	var x slot.Pos
	for x = 1; x <= ReelNumBon[es-1]; x++ {
		if es != scat {
			g.SetSym(x, 1, es)
			g.SetSym(x, 2, es)
			g.SetSym(x, 3, es)
		} else {
			g.SetSym(x, 1, 0)
			g.SetSym(x, 2, scat)
			g.SetSym(x, 3, 0)
		}
	}
	var s = slot.NewSnapFlat(sn, 5)

	var calc = func(w io.Writer) (float64, float64) {
		var N = s.Count()
		var cost = g.Cost()
		var rtpsym = s.SumPays() / cost / N
		var HRfg, EK, _ = slot.Efs(s)
		var q = EK
		var sq = 1 / (1 - q)
		if q > 0 {
			fmt.Fprintf(w, "symbols: µ = %.8g%%\n", rtpsym*100)
			fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/HRfg, q, sq)
		}
		fmt.Fprintf(w, "RTP[%d] = %.8g%%\n", es, rtpsym*100)
		fmt.Fprintln(w)
		return rtpsym, math.NaN()
	}

	func() {
		var ctx2, cancel2 = context.WithCancel(ctx)
		defer cancel2()
		var reels = g.GetReels(sp.MRTP)
		switch ReelNumBon[g.ES-1] {
		case 2:
			sp.Total = uint64(len(reels.Reel(3))-2) * uint64(len(reels.Reel(4))-2) * uint64(len(reels.Reel(5))-2)
			BruteForce5x3es2(ctx2, sp, s, g, reels)
		case 3:
			sp.Total = uint64(len(reels.Reel(4))-2) * uint64(len(reels.Reel(5))-2)
			BruteForce5x3es3(ctx2, sp, s, g, reels)
		default:
			panic("no way to here")
		}
	}()
	var N = s.Count()
	var rtp, sigma = calc(os.Stdout)
	var q = float64(s.A[0].Cfs) / N
	return rtp, sigma, q
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("\n(1/2) bonus reels calculations\n")
	var rtpe, qe [sn]float64
	var es slot.Sym
	for es = 2; es <= scat; es++ {
		fmt.Printf("*calculations for expanding symbol [%d]*\n", es)
		rtpe[es-1], _, qe[es-1] = CalcStatBon(ctx, sp, es)
		if ctx.Err() != nil {
			return 0, 0
		}
	}
	var rtpsym, qfs float64
	for _, es := range ReelExpSym {
		rtpsym += rtpe[es-1]
		qfs += qe[es-1]
	}
	rtpsym /= float64(len(ReelExpSym))
	qfs /= float64(len(ReelExpSym))
	var sqfs = 1 / (1 - qfs)
	var rtpfs = sqfs * rtpsym
	fmt.Printf("free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 10/qfs, qfs, sqfs)
	fmt.Printf("RTPfs = sq*rtp(sym) = %.5g*%.5g = %.8g%%\n", sqfs, rtpsym*100, rtpfs*100)

	fmt.Printf("\n(2/2) regular reels calculations\n")
	var g = NewGame(sp.Sel)
	var s = slot.NewSnapFlat(sn, 5)

	var calc = func(w io.Writer) (float64, float64) {
		var N = s.Count()
		var cost = g.Cost()
		var rtpsym = s.SumPays() / cost / N
		var HRfg, EK, _ = slot.Efs(s)
		var q = EK
		var sq = 1 / (1 - q)
		var rtp = rtpsym + q*rtpfs
		fmt.Fprintf(w, "symbols: µ = %.8g%%\n", rtpsym*100)
		fmt.Fprintf(w, "free: HRfg = 1/%.5g, q = %.5g, sq = 1/(1-q) = %.5g\n", 1/HRfg, q, sq)
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g*%.5g(fg) = %.8g%%\n", rtpsym*100, q, rtpfs*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
