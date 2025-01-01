package printer

import (
	"fmt"

	"github.com/mstergianis/pacdiff/pkg/diff"
)

func NewPrinter(options ...Option) *Printer {
	p := &Printer{}
	WithDepth("  ")(p)
	for _, opt := range options {
		opt(p)
	}
	return p
}

type Option func(*Printer)

type Printer struct {
	depthMarker string
}

func WithDepth(depthMarker string) Option {
	return func(p *Printer) {
		p.depthMarker = depthMarker
	}
}

func (p *Printer) Print(d diff.Diff) {
	p.print(0, d)
}

func (p *Printer) print(depth int, d diff.Diff) {
	for k, v := range d {
		p.printDepthMarker(depth)
		fmt.Printf("%q: ", k)
		switch vConcrete := v.(type) {
		case map[string]any:
			{
				fmt.Printf("{")
				fmt.Printf("\n")
				p.print(depth+1, diff.Diff(vConcrete))
				p.printDepthMarker(depth)
				fmt.Printf("}")
			}
		case string:
			{
				fmt.Printf("%q", vConcrete)
			}
		default:
			panic("unimplemented")
		}
		fmt.Printf("\n")
	}
}

func (p *Printer) printDepthMarker(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print(p.depthMarker)
	}
}

func (p *Printer) PrintUnified(d diff.Diff) {
	panic("Printer.PrintUnified - unimplemented")
}
