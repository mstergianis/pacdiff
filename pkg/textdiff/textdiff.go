package textdiff

import (
	"bytes"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/mstergianis/pacdiff/pkg/diff"
)

func Myer(left, leftName, right, rightName string) (string, error) {
	lLines := strings.Split(strings.TrimSuffix(left, "\n"), "\n")
	rLines := strings.Split(strings.TrimSuffix(right, "\n"), "\n")

	n, m := len(lLines), len(rLines)
	maxMoves := n + m
	trace := shortestEditTrace(lLines, rLines, n, m, maxMoves)
	coords := backtrack(trace, n, m, maxMoves)

	diffs := []diff.Diff{}
	for c1, c2 := range getCoordPairs(values(slices.Backward(coords))) {
		if c1.X < c2.X && c1.Y < c2.Y {
			diffs = append(diffs, diff.Diff{
				Typ:     diff.Equality,
				Content: lLines[c1.X],
			})
		} else if c1.X < c2.X {
			diffs = append(diffs, diff.Diff{
				Typ:     diff.Deletion,
				Content: lLines[c1.X],
			})
		} else {
			diffs = append(diffs, diff.Diff{
				Typ:     diff.Insertion,
				Content: rLines[c1.Y],
			})
		}
	}

	b := &bytes.Buffer{}
	for _, diff := range diffs {
		fmt.Fprintln(b, diff)
	}

	return b.String(), nil
}

func getCoordPairs(seq iter.Seq[Coord]) iter.Seq2[Coord, Coord] {
	return func(yield func(c1, c2 Coord) bool) {
		next, stop := iter.Pull(seq)
		defer stop()

		var (
			v1  Coord
			ok1 bool
		)
		v1, ok1 = next()
		if !ok1 {
			return
		}
		for {
			v2, ok2 := next()
			if !ok2 {
				return
			}
			if !yield(v1, v2) {
				return
			}
			v1 = v2
		}
	}
}

func values(seq2 iter.Seq2[int, Coord]) iter.Seq[Coord] {
	return func(yield func(c Coord) bool) {
		for _, v := range seq2 {
			if !yield(v) {
				return
			}
		}
	}
}

func backtrack(trace [][]int, x, y, maxMoves int) []Coord {
	coords := make([]Coord, 0, maxMoves)
	for d, v := range slices.Backward(trace) {
		k := x - y

		var prevK int
		if k == -d || (k != d && v[k-1+maxMoves] < v[k+1+maxMoves]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		var (
			prevX = v[prevK+maxMoves]
			prevY = prevX - prevK
		)
		for x >= prevX && y >= prevY {
			coords = append(coords, Coord{X: x, Y: y})
			x, y = x-1, y-1
		}
		x, y = prevX, prevY
	}
	return coords
}

func shortestEditTrace(lLines, rLines []string, n, m, maxMoves int) (trace [][]int) {
	trace = make([][]int, 0, maxMoves)
	for d := 0; d <= maxMoves; d++ {
		v := make([]int, 2*maxMoves+1)
		if d > 0 {
			copy(v, trace[d-1])
		}
		trace = append(trace, v)
		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[k-1+maxMoves] < v[k+1+maxMoves]) {
				x = v[k+1+maxMoves]
			} else {
				x = v[k-1+maxMoves] + 1
			}
			y := x - k
			for x < n && y < m && lLines[x] == rLines[y] {
				x++
				y++
			}
			v[k+maxMoves] = x
			if x >= n && y >= m {
				return
			}
		}
	}

	return trace
}

type Coord struct {
	X int
	Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}
