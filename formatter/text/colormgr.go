package text

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

const escSeqFmt = "\033[%dm"

type colorMgr struct {
	colors    [iface.LevelCount]Color
	markColor Color

	colorSeqs    [iface.LevelCount][]byte
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

func (mgr *colorMgr) Color(level iface.Level) Color {
	return mgr.colors[level]
}

func (mgr *colorMgr) SetColor(level iface.Level, color Color) {
	mgr.colors[level] = color
	mgr.colorSeqs[level] = makeColorSeq(color)
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
	mgr.markColorSeq = makeColorSeq(color)
}

func (mgr *colorMgr) ColorEars(level iface.Level) ([]byte, []byte) {
	return mgr.colorSeqs[level], mgr.resetSeq
}

func (mgr *colorMgr) MarkColorEars() ([]byte, []byte) {
	return mgr.markColorSeq, mgr.resetSeq
}

func makeColorSeq(color Color) []byte {
	return []byte(fmt.Sprintf(escSeqFmt, color))
}
