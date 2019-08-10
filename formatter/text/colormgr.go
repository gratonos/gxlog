package text

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

const escSeqFmt = "\033[%dm"

type colorMgr struct {
	colors    [iface.LogLevelCount]Color
	markColor Color

	colorSeqs    [iface.LogLevelCount][]byte
	markColorSeq []byte
	resetSeq     []byte
}

func newColorMgr(colorMap map[iface.Level]Color, markColor Color) *colorMgr {
	mgr := &colorMgr{
		resetSeq: makeColorSeq(0),
	}
	mgr.MapColors(colorMap)
	mgr.SetMarkColor(markColor)
	return mgr
}

func (this *colorMgr) Color(level iface.Level) Color {
	return this.colors[level]
}

func (this *colorMgr) SetColor(level iface.Level, color Color) {
	this.colors[level] = color
	this.colorSeqs[level] = makeColorSeq(color)
}

func (this *colorMgr) MapColors(colorMap map[iface.Level]Color) {
	for level, color := range colorMap {
		this.SetColor(level, color)
	}
}

func (this *colorMgr) MarkColor() Color {
	return this.markColor
}

func (this *colorMgr) SetMarkColor(color Color) {
	this.markColor = color
	this.markColorSeq = makeColorSeq(color)
}

func (this *colorMgr) ColorEars(level iface.Level) ([]byte, []byte) {
	return this.colorSeqs[level], this.resetSeq
}

func (this *colorMgr) MarkColorEars() ([]byte, []byte) {
	return this.markColorSeq, this.resetSeq
}

func makeColorSeq(color Color) []byte {
	return []byte(fmt.Sprintf(escSeqFmt, color))
}
