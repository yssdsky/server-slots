package slot

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"sync"
	"time"
)

// Function to report about progress of calculation by brute force
func ReporterBF(ctx context.Context, snapch []chan Simulator2,
	sp *ScanPar, ss Simulator2, g SlotGeneric, calc CalcFunc) {
	var t0 = time.Now()
	const stepdur = 1000 * time.Millisecond
	var steps = time.Tick(stepdur)
	if sel := g.GetSel(); sel > 0 {
		fmt.Printf("selected %d lines\n", sel)
	}
	fmt.Printf("calculation started...\r")
	var total = float64(sp.Total)
	var (
		dur time.Duration
		N   float64
	)
	var latest = make([]Simulator2, len(snapch))
	var closed = make([]bool, len(snapch))
	var activeCount = len(snapch)
	for activeCount > 0 {
		select {
		case <-ctx.Done():
			activeCount = 0
		case <-steps:
		}
		for i, ch := range snapch {
			if !closed[i] {
				select {
				case snap, ok := <-ch:
					if !ok {
						closed[i] = true
						activeCount--
					} else {
						latest[i] = snap
					}
				default:
				}
			}
		}

		ss.Clear()
		for i := range latest {
			if latest[i] != nil {
				ss.Add(latest[i])
			}
		}

		dur = time.Since(t0)
		N = ss.Count()
		var RTP, _ = calc(io.Discard)

		var exp string
		if dur > 2*stepdur {
			exp = time.Duration(float64(dur-stepdur) * total / N).Truncate(time.Second).String()
		} else {
			exp = "?"
		}
		fmt.Printf("processed %.1fm/%.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%  \r",
			N/1e6, total/1e6, N/total*100, dur.Truncate(time.Second), exp, RTP*100)
	}

	// report results
	fmt.Printf("completed %.5g%% (%d), time spent %v                    \n",
		N/total*100, int(N), dur)
	var reels = g.GetReels(sp.MRTP)
	fmt.Println()
	fmt.Printf("reels lengths %s, total reshuffles %d\n", reels.String(0), reels.Reshuffles(0))
	calc(os.Stdout)
}

func WorkerBF(ctx context.Context, outch chan Simulator2, ti, tn uint64,
	sp *ScanPar, local Simulator2, wgame SlotGeneric) {
	// Grid dimension.
	var sx, sy = wgame.Dim()
	var reels = wgame.GetReels(sp.MRTP).Clone().SetPadding(sy)
	var total = reels.Reshuffles(sy)
	var _, iscascade = wgame.(SlotCascade)

	// Precompute reel lengths.
	var rlen = make([]int, sx)
	for x := range sx {
		rlen[x] = len(reels[x]) - int(sy)
	}
	// Precompute delimiters for position calculation.
	var delim = make([]uint64, sx)
	delim[0] = 1
	var x Pos
	for x = 1; x < sx; x++ {
		delim[x] = delim[x-1] * uint64(rlen[x-1])
	}
	var wins Wins
	var rpos = make([]int, sx)
	for x := range sx {
		rpos[x] = -1
	}
	var i, N uint64
	var pos int
	// Using one general loop instead of five nested loops
	// gives ~30% performance increase.
	for i = ti; i < total; i += tn {
		N++
		if N&0xFFF == 0 {
			select {
			case <-ctx.Done():
				return
			case outch <- local.Clone():
			default:
			}
		}
		for x = range sx {
			pos = int(i/delim[x]) % rlen[x]
			if rpos[x] != pos {
				wgame.SetColFast(x+1, reels[x], pos)
				rpos[x] = pos
			} else {
				break
			}
		}
		local.Simulate(wgame, &wins, sp.MRTP)
		if iscascade && len(wins) > 0 {
			for x := range sx {
				rpos[x] = -1
			}
		}
		wins.Reset()
	}
	// clear channel buffer
	select {
	case <-outch:
	default:
	}
	// push final result to channel
	outch <- local
}

func BruteForcex2(ctx context.Context,
	sp *ScanPar, ss Simulator2, mgame SlotGeneric, calc CalcFunc) {
	var reels = mgame.GetReels(sp.MRTP)
	if sp.TN%len(reels[0]) == 0 {
		panic("BruteForcex: thread number equals to 1-st reel length")
	}

	var snapch = make([]chan Simulator2, sp.TN)
	for i := range snapch {
		snapch[i] = make(chan Simulator2, 1)
	}

	var wg1, wg2 sync.WaitGroup
	func() {
		if calc != nil {
			var ctxrep, cancel = context.WithCancel(ctx)
			defer cancel()
			wg2.Go(func() {
				ReporterBF(ctxrep, snapch, sp, ss, mgame, calc)
			})
		}
		for ti, ch := range snapch {
			var local = ss.Clone()
			var wgame = mgame.Clone()
			wg1.Go(func() {
				defer close(ch)
				WorkerBF(ctx, ch, uint64(ti), uint64(len(snapch)), sp, local, wgame)
			})
		}
		wg1.Wait()
	}()
	wg2.Wait()
}

// Function to report about progress of calculation by brute force
func ProgressBF(ctx context.Context, sp *ScanPar, s Simulator, calc func(io.Writer) (float64, float64)) {
	const stepdur = 1000 * time.Millisecond
	var t0 = time.Now()
	var steps = time.Tick(stepdur)
	fmt.Printf("calculation started...\r")
	var total = float64(sp.Total)
	var (
		dur time.Duration
		N   float64
		RTP float64
	)
	var param = func() {
		dur = time.Since(t0)
		N = s.Count()
		RTP, _ = calc(io.Discard)
	}
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-steps:
			param()
			var exp = time.Duration(float64(dur) * total / N)
			fmt.Printf("processed %.1fm/%.1fm, ready %2.2f%% (%v / %v), RTP = %2.2f%%  \r",
				N/1e6, total/1e6, N/total*100,
				dur.Truncate(stepdur), exp.Truncate(stepdur),
				RTP*100)
		}
	}

	// report results
	param()
	fmt.Printf("completed %.5g%% (%d), time spent %v                    \n",
		N/total*100, int(N), dur)
}

func BruteForcex(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric) {
	// Grid dimension.
	var sx, sy = g.Dim()
	var reels = g.GetReels(sp.MRTP).Clone().SetPadding(sy)
	var total = reels.Reshuffles(sy)
	var tn64 = uint64(sp.TN)
	if sp.TN%(len(reels[0])-int(sy)) == 0 {
		panic("BruteForcex: thread number equals to 1-st reel length")
	}
	// Precompute reel lengths.
	var rlen = make([]int, sx)
	var x Pos
	for x = range sx {
		rlen[x] = len(reels[x]) - int(sy)
	}
	// Precompute delimiters for position calculation.
	var delim = make([]uint64, sx)
	delim[0] = 1
	for x = 1; x < sx; x++ {
		delim[x] = delim[x-1] * uint64(rlen[x-1])
	}

	var wg sync.WaitGroup
	wg.Add(sp.TN)
	for ti := range tn64 {
		var gt = g.Clone()
		var _, iscascade = gt.(SlotCascade)
		go func() {
			defer wg.Done()
			var wins Wins
			var rpos = make([]int, sx)
			for x := range sx {
				rpos[x] = -1
			}
			var i, N uint64
			var x Pos
			var pos int
			// Using one general loop instead of five nested loops
			// gives ~30% performance increase.
			for i = ti; i < total; i += tn64 {
				N++
				if N&0xFFF == 0 {
					select {
					case <-ctx.Done():
						return
					default:
					}
				}
				for x = range sx {
					pos = int(i/delim[x]) % rlen[x]
					if rpos[x] != pos {
						gt.SetColFast(x+1, reels[x], pos)
						rpos[x] = pos
					} else {
						break
					}
				}
				s.Simulate(gt, &wins, sp.MRTP)
				if iscascade && len(wins) > 0 {
					for x = range sx {
						rpos[x] = -1
					}
				}
				wins.Reset()
			}
		}()
	}
	wg.Wait()
}

func BruteForce5x3Big(ctx context.Context, sp *ScanPar, s Simulator, g SlotGeneric, r1, rb, r5 []Sym) {
	// var total = uint64(len(r1)) * uint64(len(rb)) * uint64(len(r5))
	r1 = append(slices.Clone(r1), r1[:3]...)
	rb = append(slices.Clone(rb), rb[:3]...)
	r5 = append(slices.Clone(r5), r5[:3]...)
	var tn64 = uint64(sp.TN)
	var wg sync.WaitGroup
	wg.Add(sp.TN)
	for ti := range tn64 {
		var gt = g.Clone()
		var cb = gt.(Bigger)
		var N uint64
		go func() {
			defer wg.Done()

			var wins Wins
			for i1 := range len(r1) - 3 {
				gt.SetColFast(1, r1, i1)
				for _, big := range rb[:len(rb)-3] {
					cb.SetBig(big)
					for i5 := range len(r5) - 3 {
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
						gt.SetColFast(5, r5, i5)
						s.Simulate(gt, &wins, sp.MRTP)
						wins.Reset()
					}
				}
			}
		}()
	}
	wg.Wait()
}
