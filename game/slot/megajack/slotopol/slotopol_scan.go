package slotopol

import (
	"context"
	"fmt"
	"io"
	"math"

	"github.com/slotopol/server/game/slot"
)

var Emje float64 // Eldorado game 1 spin calculated expectation

// Eldorado expectation.
func ExpEldorado() {
	var sum float64
	for _, v := range Eldorado {
		sum += v
	}
	Emje = sum / float64(len(Eldorado))
	fmt.Printf("eldorado 1 spin: count = %d, E = %g\n", len(Eldorado), Emje)
}

var Emjm float64 // Monopoly game calculated expectation

func ExpMonopoly() {
	var dices = [7]int{1, 1, 1, 1, 1, 1, 1}
	var sumi, sumii, sumj float64
	var count, zcount int = 6 * 6 * 6 * 6 * 6 * 6 * 6, 0
	var pos int
	var win float64
	for i1 := 1; i1 <= 6; i1++ {
		for i2 := 1; i2 <= 6; i2++ {
			for i3 := 1; i3 <= 6; i3++ {
				for i4 := 1; i4 <= 6; i4++ {
					for i5 := 1; i5 <= 6; i5++ {
						for i6 := 1; i6 <= 6; i6++ {
							for i7 := 1; i7 <= 6; i7++ {
								pos, sumj = 1, 0
								for j := 0; j < 7; j++ {
									pos = (pos+dices[j]-1)%20 + 1
									if Monopoly[pos-1].Jump > 0 {
										win = Monopoly[Monopoly[pos-1].Jump-1].Mult
									} else {
										win = Monopoly[pos-1].Mult
									}
									if Monopoly[pos-1].Dice {
										win *= float64(dices[j])
									}
									if Monopoly[pos-1].Jump > 0 {
										pos = Monopoly[pos-1].Jump
									}
									sumj += win
								}
								if sumj == 0 {
									sumj += 5000
									zcount++
								}
								sumi += sumj
								sumii += sumj * sumj

								dices[6] = dices[6]%6 + 1
							}
							dices[5] = dices[5]%6 + 1
						}
						dices[4] = dices[4]%6 + 1
					}
					dices[3] = dices[3]%6 + 1
				}
				dices[2] = dices[2]%6 + 1
			}
			dices[1] = dices[1]%6 + 1
		}
		dices[0] = dices[0]%6 + 1
	}
	Emjm = sumi / float64(count)
	var v = sumii/float64(count) - Emjm*Emjm
	var sigma = math.Sqrt(v)
	fmt.Printf("monopoly: count = %d, sum = %g, zerocount = %d, p(zero) = 1/%d, E = %g\n", count, sumi, zcount, int(float64(count)/float64(zcount)), Emjm)
	fmt.Printf("monopoly: variance = %.6g, sigma = %.6g, limits = %.6g ... %.6g\n", v, sigma, Emjm-sigma, Emjm+sigma)
}

func CalcStat(ctx context.Context, sp *slot.ScanPar) (float64, float64) {
	fmt.Printf("*bonus games calculations*\n")
	ExpEldorado()
	ExpMonopoly()
	fmt.Printf("*reels calculations*\n")
	var g = NewGame(sp.Sel)
	var s = slot.NewStatSnap(1, sn, 5, mjap, mjj)

	var calc = func(w io.Writer) (float64, float64) {
		var N = s.Count()
		var cost = g.Cost()
		var rtpsym = s.SumPays() / cost / N
		var HRmje9 = float64(s.A[0].Cb[mje9-1]) / N
		var qmje9 = HRmje9 / float64(g.Sel)
		var rtpmje9 = Emje * 9 * qmje9
		var HRmjm = float64(s.A[0].Cb[mjm-1]) / N
		var qmjm = HRmjm / float64(g.Sel)
		var rtpmjm = Emjm * qmjm
		var rtp = rtpsym + rtpmje9 + rtpmjm
		fmt.Fprintf(w, "symbols: µ = %.6f%%\n", rtpsym*100)
		fmt.Fprintf(w, "spin9 bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", 1/HRmje9, rtpmje9*100)
		fmt.Fprintf(w, "monopoly bonuses: hit rate 1/%.5g, rtp = %.6f%%\n", 1/HRmjm, rtpmjm*100)
		if jn := s.A[0].Cj[mjj-1]; jn > 0 {
			var HRj = float64(jn) / N
			fmt.Fprintf(w, "jackpots: count %d, hit rate 1/%.12g\n", jn, 1/HRj)
		}
		fmt.Fprintf(w, "RTP = %.5g(sym) + %.5g(mje9) + %.5g(mjm) = %.6f%%\n", rtpsym*100, rtpmje9*100, rtpmjm*100, rtp*100)
		return rtp, math.NaN()
	}

	return slot.ScanReelsCommon2(ctx, sp, s, g, calc)
}
