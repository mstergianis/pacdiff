package diff

import (
	"fmt"
)

type Hunk struct {
	LeftName string

	// LeftStart as a line number
	LeftStart int
	// LeftEnd as a line number
	LeftEnd int

	RightName string

	// RightStart as a line number
	RightStart int
	// RightEnd as a line number
	RightEnd int

	Diffs []Diff
}

func (h Hunk) String() string {
	s := fmt.Sprintf("@@ -%s +%s @@\n",
		fmtHunkLines(h.LeftStart, h.LeftEnd),
		fmtHunkLines(h.RightStart, h.RightEnd),
	)

	for _, d := range h.Diffs {
		s += d.String() + "\n"
	}

	return s
}

func fmtHunkLines(start, end int) string {
	if start == end {
		return fmt.Sprintf("%d", start)
	}

	return fmt.Sprintf("%d,%d", start, end)
}

type Diff struct {
	Typ     DiffTyp
	Content string
}

func (d Diff) String() string {
	return d.Typ.String() + d.Content
}

type DiffTyp int

func (d DiffTyp) String() string {
	switch d {
	case Insertion:
		return "+"
	case Deletion:
		return "-"
	case Equality:
		return " "
	}
	panic(fmt.Sprintf("error: encountered an unknown diff.DiffTyp %d", d))
}

const (
	Insertion = iota
	Deletion
	Equality
)
