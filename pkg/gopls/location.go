package main

import (
	"fmt"
	"path/filepath"

	"github.com/charlievieth/xtools/span"
)

type Location struct {
	Path  string `json:"path"`
	Point Point  `json:"point"`
}

func NewLocation(path string, point Point) Location {
	return Location{
		Path:  filepath.Clean(path),
		Point: point,
	}
}

func (l Location) URI() span.URI     { return span.URIFromPath(l.Path) }
func (l Location) Start() span.Point { return l.Point.Point() }
func (l Location) End() span.Point   { return l.Start() }
func (l Location) String() string    { return fmt.Sprintf("%s", l.Span()) }

func (l Location) Span() span.Span {
	p := l.Start()
	return span.New(l.URI(), p, p)
}

type Point struct {
	Line   int `json:"line"`
	Column int `json:"column"`
	Offset int `json:"offset"`
}

func NewPoint(line, column, offset int) Point {
	return Point{
		Line:   line,
		Column: column,
		Offset: offset,
	}
}

func (p Point) Point() span.Point {
	return span.NewPoint(p.Line, p.Column, p.Offset)
}
