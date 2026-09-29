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

// Measures classical Dijkstra on the same 16x16 game-map grid used by the race-logic model, so any machine
// (Apple M-series, AMD, Intel, ...) can contribute a real baseline to the comparison.
func main() {
	side := flag.Int("side", 16, "grid side (side x side map)")
	queries := flag.Int("queries", 20000, "timed single-source queries")
	flag.Parse()

	g := optical.NewGridGraph(*side, *side, 15, 42)
	n := len(g.Adj)
	for i := 0; i < 2000; i++ { // warm-up
		optical.Dijkstra(g, i%n)
	}
	start := time.Now()
	for i := 0; i < *queries; i++ {
		optical.Dijkstra(g, i%n)
	}
	perQuery := time.Since(start) / time.Duration(*queries)

	fmt.Printf("CPU:        %s\n", cpuModel())
	fmt.Printf("Platform:   %s/%s, Go %s, single thread\n", runtime.GOOS, runtime.GOARCH, runtime.Version())
	fmt.Printf("Map:        %dx%d (%d nodes, %d directed edges)\n", *side, *side, n, g.NumEdges())
	fmt.Printf("Dijkstra:   %.2f us per single-source query (%d queries)\n", float64(perQuery.Nanoseconds())/1000, *queries)
	fmt.Printf("SilicaCore: 42.2 ns per query (simulated, 16x16) -> %.0fx\n", float64(perQuery.Nanoseconds())/42.2)
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
