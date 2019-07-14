package text

import (
	"github.com/gratonos/gxlog/iface"
)

type Color int

const (
	FullHeader = "{{time}} {{level}} {{file}}:{{line}} {{pkg}}.{{func}} " +
		"{{prefix}}[{{context}}] {{msg}}\n"
	CompactHeader = "{{time:time.us}} {{level}} {{file:1}}:{{line}} " +
		"{{pkg}}.{{func}} {{prefix}}[{{context}}] {{msg}}\n"
	SyslogHeader = "{{file:1}}:{{line}} {{pkg}}.{{func}} {{prefix}}[{{context}}] {{msg}}\n"
)

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

type Config struct {
	// The Header is the format specifier of a text formatter.
	// It is used to specify which and how the fields of a Record to be formatted.
	// The pattern of a field specifier is {{<name>[:property][%fmtstr]}}.
	// e.g. {{level:char}}, {{line%05d}}, {{pkg:1}}, {{context:list%40s}}
	// All fields have support for the fmtstr. If the fmtstr is NOT the default one
	// of a field, it will be passed to fmt.Sprintf to format the field and this
	// affects the performance a little.
	// The supported properties vary with fields.
	// All supported fields are as the follows:
	//    name    | supported property       | defaults     | property examples
	//  ----------+--------------------------+--------------+------------------------
	//    time    | <date|time>[.ms|.us|.ns] | "date.us" %s | "date.ns", "time"
	//            | layout that is supported |              | time.RFC3339Nano
	//            |   by the time package    |              | "02 Jan 06 15:04 -0700"
	//    level   | <full|char>              | "full"    %s | "full", "char"
	//    file    | <lastSegs>               | 0         %s | 0, 1, 2, ...
	//    line    |                          |           %d |
	//    pkg     | <lastSegs>               | 0         %s | 0, 1, 2, ...
	//    func    | <lastSegs>               | 0         %s | 0, 1, 2, ...
	//    prefix  |                          |           %s |
	//    context | <pair|list>              | "pair"    %s | "pair", "list"
	//    msg     |                          |           %s |
	Header    string
	Coloring  bool
	ColorMap  map[iface.Level]Color
	MarkColor Color
}

func (this *Config) SetDefaults() {
	if this.Header == "" {
		this.Header = CompactHeader
	}

	colorMap := map[iface.Level]Color{
		iface.Trace: Green,
		iface.Debug: Green,
		iface.Info:  Green,
		iface.Warn:  Yellow,
		iface.Error: Red,
		iface.Fatal: Red,
	}
	for level, color := range this.ColorMap {
		colorMap[level] = color
	}
	this.ColorMap = colorMap

	if this.MarkColor == 0 {
		this.MarkColor = Magenta
	}
}
