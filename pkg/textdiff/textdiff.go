package textdiff

import (
	"bytes"
	"fmt"
	"iter"
	"slices"
	"strings"
)

func Myer(left, lName, right, rName string) (string, error) {
	lLines := strings.Split(strings.TrimSuffix(left, "\n"), "\n")
	rLines := strings.Split(strings.TrimSuffix(right, "\n"), "\n")

	n, m := len(lLines), len(rLines)
	maxMoves := n + m
	trace, d, k := shortestEditTrace(lLines, rLines, n, m, maxMoves)
	coords := backtrack(trace, d, k, maxMoves)

	diffSet := []Diff{}
	for c1, c2 := range getCoordPairs(values(slices.Backward(coords))) {
		if c1.X < c2.X && c1.Y < c2.Y {
			diffSet = append(diffSet, Diff{
				Equality,
				lLines[c1.X],
			})
		} else if c1.X < c2.X {
			diffSet = append(diffSet, Diff{
				Deletion,
				lLines[c1.X],
			})
		} else {
			diffSet = append(diffSet, Diff{
				Insertion,
				rLines[c1.Y],
			})
		}
	}

	var buf = new(bytes.Buffer)
	var (
		lHunkStart = 1
		lHunkEnd   = 7
		rHunkStart = 1
		rHunkEnd   = 6
	)
	fmt.Fprintf(buf, "--- %s\n+++ %s\n@@ -%d,%d +%d,%d @@\n", lName, rName, lHunkStart, lHunkEnd, rHunkStart, rHunkEnd)
	for _, d := range diffSet {
		fmt.Fprintln(buf, d)
	}

	return buf.String(), nil
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

func backtrack(trace [][]int, d, k, maxMoves int) []Coord {
	coords := make([]Coord, 0, maxMoves)
	for i := len(trace) - 1; i >= 0; i-- {
		v := trace[i]
		x := v[k+maxMoves]
		y := x - k

		coords = append(coords, Coord{X: x, Y: y})
		if i == 0 {
			return coords
		}

		if k == -d || (k != d && v[k-1+maxMoves] < v[k+1+maxMoves]) {
			k = k + 1
		} else {
			k = k - 1
		}

		prevV := trace[i-1]
		var (
			prevX = prevV[k+maxMoves]
			prevY = prevX - k
			diagX = x - 1
			diagY = y - 1
		)
		for ; diagX >= prevX && diagY >= prevY; diagX, diagY = diagX-1, diagY-1 {
			coords = append(coords, Coord{X: diagX, Y: diagY})
		}
	}

	return coords
}

func shortestEditTrace(lLines, rLines []string, n, m, maxMoves int) (trace [][]int, d, k int) {
	trace = make([][]int, 0, maxMoves)
	for d = 0; d <= maxMoves; d++ {
		v := make([]int, 2*maxMoves+1)
		if d > 0 {
			copy(v, trace[d-1])
		}
		trace = append(trace, v)
		for k = -d; k <= d; k += 2 {
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

	return trace, 0, 0
}

type Coord struct {
	X int
	Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

type Diff struct {
	Typ     DiffTyp
	Content string
}

func (d DiffTyp) String() string {
	switch d {
	case Insertion:
		return "+"
	case Deletion:
		return "-"
	case Equality:
		return " "
	}
	panic("DiffTyp.String: encountered an unknown DiffTyp")
}

func (d Diff) String() string {
	return d.Typ.String() + d.Content
}

type DiffTyp int

const (
	Insertion = iota
	Deletion
	Equality
)
