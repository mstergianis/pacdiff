package printer

import (
	"fmt"
	"io"
	"os"

	"github.com/mstergianis/pacdiff/pkg/differ"
)

func NewPrinter(options ...Option) *Printer {
	p := &Printer{}
	WithOutputWriter(os.Stdout)
	WithErrorWriter(os.Stderr)
	WithDepth("  ")(p)
	for _, opt := range options {
		opt(p)
	}
	return p
}

type Option func(*Printer)

type Printer struct {
	depthMarker  string
	outputWriter io.Writer
	errorWriter  io.Writer
}

func WithDepth(depthMarker string) Option {
	return func(p *Printer) {
		p.depthMarker = depthMarker
	}
}

func WithOutputWriter(outputWriter io.Writer) Option {
	return func(p *Printer) {
		p.outputWriter = outputWriter
	}
}

func WithErrorWriter(errorWriter io.Writer) Option {
	return func(p *Printer) {
		p.errorWriter = errorWriter
	}
}

func (p *Printer) printDepthMarker(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Print(p.depthMarker)
	}
}

func (p *Printer) PrintUnified(fileGroupedHunks differ.GroupedHunksSlice) {
	for _, fileGroup := range fileGroupedHunks {
		p.printUnifiedHeader(fileGroup.LeftFile, fileGroup.RightFile)

		for _, hunk := range fileGroup.Hunks {
			p.print(hunk)
		}
	}
}

func (p *Printer) printf(format string, args ...interface{}) {
	fmt.Fprintf(p.outputWriter, format, args...)
}

func (p *Printer) print(a ...any) {
	fmt.Fprint(p.outputWriter, a...)
}

func (p *Printer) println(a ...any) {
	fmt.Fprintln(p.outputWriter, a...)
}

func (p *Printer) errorf(format string, args ...interface{}) {
	fmt.Fprintf(p.errorWriter, format, args...)
}

func (p *Printer) printUnifiedHeader(left, right string) {
	p.printf("--- %s\n", left)
	p.printf("+++ %s\n", right)
}
