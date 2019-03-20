package text

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

type Color int

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

const escSeqFmt = "\033[%dm"

type colorMgr struct {
	colors    []Color
	markColor Color

	colorSeqs    [][]byte
	markColorSeq []byte
	resetSeq     []byte
}

func newColorMgr() *colorMgr {
	colors := []Color{
		iface.Trace: Green,
		iface.Debug: Green,
		iface.Info:  Green,
		iface.Warn:  Yellow,
		iface.Error: Red,
		iface.Fatal: Red,
	}
	mgr := &colorMgr{
		colors:       colors,
		markColor:    Magenta,
		colorSeqs:    initColorSeqs(colors),
		markColorSeq: makeSeq(Magenta),
		resetSeq:     makeSeq(0),
	}
	return mgr
}

func (mgr *colorMgr) Color(level iface.Level) Color {
	return mgr.colors[level]
}

func (mgr *colorMgr) SetColor(level iface.Level, color Color) {
	mgr.colors[level] = color
	mgr.colorSeqs[level] = makeSeq(color)
}

func (mgr *colorMgr) MapColors(colorMap map[iface.Level]Color) {
	for level, color := range colorMap {
		mgr.SetColor(level, color)
	}
}

func (mgr *colorMgr) MarkColor() Color {
	return mgr.markColor
}

func (mgr *colorMgr) SetMarkColor(color Color) {
	mgr.markColor = color
	mgr.markColorSeq = makeSeq(color)
}

func (mgr *colorMgr) ColorEars(level iface.Level) ([]byte, []byte) {
	return mgr.colorSeqs[level], mgr.resetSeq
}

func (mgr *colorMgr) MarkColorEars() ([]byte, []byte) {
	return mgr.markColorSeq, mgr.resetSeq
}

func initColorSeqs(colors []Color) [][]byte {
	colorSeqs := make([][]byte, len(colors))
	for i := range colors {
		colorSeqs[i] = makeSeq(colors[i])
	}
	return colorSeqs
}

func makeSeq(color Color) []byte {
	return []byte(fmt.Sprintf(escSeqFmt, color))
}
