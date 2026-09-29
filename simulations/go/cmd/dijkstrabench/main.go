package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"devaneio-ricardo/tof-cpu/pkg/optical"
)

// Measures the electronic baselines for one single-source query (all distances)
// on the same grid map used by the race-logic simulation:
//   - Dijkstra with a binary heap (the baseline of the preliminary article);
//   - Dial's bucket queue, the best algorithm for integer weights 1..15;
//   - copying one row of a precomputed all-pairs table (static map).
func main() {
	side := flag.Int("side", 16, "grid side (side x side map)")
	queries := flag.Int("queries", 20000, "timed single-source queries")
	raceNS := flag.Float64("race-ns", 42.2, "simulated race-logic time per query (ns)")
	programNS := flag.Float64("program-ns", 1000, "time to reprogram the map on the chip (ns, ASSUMPTION)")
	flag.Parse()

	const maxWeight = 15
	g := optical.NewGridGraph(*side, *side, maxWeight, 42)
	n := len(g.Adj)

	heapUS := timeUS(*queries, func(i int) { optical.Dijkstra(g, i%n) })

	st := optical.NewDialState(n, maxWeight)
	dialUS := timeUS(*queries, func(i int) { optical.DijkstraDial(g, i%n, st) })

	start := time.Now()
	table := optical.AllPairsTable(g, maxWeight)
	buildMS := float64(time.Since(start).Nanoseconds()) / 1e6
	row := make([]int32, n)
	rowUS := timeUS(*queries*50, func(i int) { copy(row, table[(i*7919)%n]) })

	fmt.Printf("CPU:        %s\n", cpuModel())
	fmt.Printf("Platform:   %s/%s, Go %s, single thread\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
	fmt.Printf("Map:        %dx%d (%d nodes, %d directed edges), weights 1..%d\n\n", *side, *side, n, g.NumEdges(), maxWeight)

	fmt.Printf("%-44s %12s %14s\n", "Baseline (one source, all distances)", "time", "race speedup")
	fmt.Printf("%-44s %9.2f us %13.0fx\n", "Dijkstra, binary heap", heapUS, heapUS*1000 / *raceNS)
	fmt.Printf("%-44s %9.2f us %13.0fx\n", "Dial, bucket queue", dialUS, dialUS*1000 / *raceNS)
	fmt.Printf("%-44s %9.3f us %13.1fx\n", "All-pairs table, copy one row (static map)", rowUS, rowUS*1000 / *raceNS)
	fmt.Printf("  table build: %.2f ms once, %d KB\n\n", buildMS, n*n*4/1024)

	fmt.Printf("Map changes between queries (chip reprograms in %.0f ns, ASSUMPTION):\n", *programNS)
	fmt.Printf("%-24s %14s %14s %14s\n", "queries per map change", "chip", "CPU (best)", "chip speedup")
	for _, k := range []int{1, 10, 100, 1000, 100000} {
		chipUS := (*programNS + float64(k)*(*raceNS)) / 1000
		dialTotal := float64(k) * dialUS
		tableTotal := buildMS*1000 + float64(k)*rowUS
		best := dialTotal
		if tableTotal < best {
			best = tableTotal
		}
		fmt.Printf("%-24d %11.2f us %11.2f us %13.1fx\n", k, chipUS, best, best/chipUS)
	}
}

func timeUS(count int, fn func(i int)) float64 {
	for i := 0; i < count/10+1; i++ { // warm-up
		fn(i)
	}
	start := time.Now()
	for i := 0; i < count; i++ {
		fn(i)
	}
	return float64(time.Since(start).Nanoseconds()) / float64(count) / 1000
}

func cpuModel() string {
	switch runtime.GOOS {
	case "darwin":
		if out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
			return strings.TrimSpace(string(out))
		}
	case "linux":
		if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if strings.HasPrefix(line, "model name") {
					return strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
				}
			}
		}
	}
	return "unknown"
}
