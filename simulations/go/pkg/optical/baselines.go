package optical

// Stronger electronic baselines for the shortest-path comparison.
//
// The race accelerator answers one single-source query with all distances.
// A fair CPU comparison must use the best algorithm for this workload:
//   - Dial's algorithm (bucket queue) is optimal for small integer weights
//     such as the 1..15 terrain costs used here;
//   - on a static map, an all-pairs table answers a query by copying one row.

// DialState holds reusable buffers so repeated queries do not allocate.
type DialState struct {
	dist    []int
	buckets [][]int32
}

// NewDialState prepares buffers for graphs with n nodes and weights up to maxWeight.
func NewDialState(n, maxWeight int) *DialState {
	b := make([][]int32, maxWeight+1)
	for i := range b {
		b[i] = make([]int32, 0, 64)
	}
	return &DialState{dist: make([]int, n), buckets: b}
}

// DijkstraDial returns the same distances as Dijkstra using a circular bucket
// queue of maxWeight+1 buckets (Dial, 1969). The returned slice is owned by
// the state and is overwritten by the next call.
func DijkstraDial(g Graph, src int, st *DialState) []int {
	const unreached = int(^uint(0) >> 1)
	dist := st.dist
	for i := range dist {
		dist[i] = unreached
	}
	nb := len(st.buckets)
	for i := range st.buckets {
		st.buckets[i] = st.buckets[i][:0]
	}
	dist[src] = 0
	st.buckets[0] = append(st.buckets[0], int32(src))
	pending := 1
	for d := 0; pending > 0; d++ {
		b := d % nb
		for len(st.buckets[b]) > 0 {
			last := len(st.buckets[b]) - 1
			u := int(st.buckets[b][last])
			st.buckets[b] = st.buckets[b][:last]
			pending--
			if dist[u] != d {
				continue // stale entry
			}
			for _, e := range g.Adj[u] {
				nd := d + e.Weight
				if nd < dist[e.To] {
					dist[e.To] = nd
					st.buckets[nd%nb] = append(st.buckets[nd%nb], int32(e.To))
					pending++
				}
			}
		}
	}
	for i := range dist {
		if dist[i] == unreached {
			dist[i] = -1
		}
	}
	return dist
}

// AllPairsTable precomputes every single-source distance row of a static graph.
func AllPairsTable(g Graph, maxWeight int) [][]int32 {
	n := len(g.Adj)
	st := NewDialState(n, maxWeight)
	table := make([][]int32, n)
	for s := 0; s < n; s++ {
		d := DijkstraDial(g, s, st)
		row := make([]int32, n)
		for v := range d {
			row[v] = int32(d[v])
		}
		table[s] = row
	}
	return table
}
