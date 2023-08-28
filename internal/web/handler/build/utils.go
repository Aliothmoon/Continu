package build

import "strings"

type LogWriteCloser struct {
	bid  int32
	rest []rune
	buf  []string
}

func (l *LogWriteCloser) WriteBuf(s string) {
	l.buf = append(l.buf, s)
	if len(l.buf) > 50 {
		createLog(l.bid, strings.Join(l.buf, "\n"))
		l.buf = nil
	}
}
func (l *LogWriteCloser) FlushBuf() {
	if len(l.buf) > 0 {
		createLog(l.bid, strings.Join(l.buf, "\n"))
	}
}

func (l *LogWriteCloser) Write(p []byte) (n int, err error) {
	s := string(p)
	var si = -1
	for i := range s {
		if s[i] == '\n' || s[i] == '\r' {
			var r []rune
			if si == -1 {
				si = 0
				if l.rest != nil {
					r = append(r, l.rest...)
				}
				l.rest = nil
			}
			r = append(r, []rune(s[si:i])...)

			lg := string(r)
			if lg != "" {
				l.WriteBuf(lg)
				//logger.Debug(lg)
			}

			si = i + 1
		}
	}

	if si == -1 {
		l.rest = append(l.rest, []rune(s)...)
	} else if si != len(s) {
		l.rest = append(l.rest, []rune(s[si:])...)
	}

	return len(p), nil
}

func (l *LogWriteCloser) Close() error {
	if l.rest != nil {
		lg := string(l.rest)
		if lg != "" {
			l.WriteBuf(lg)
			//logger.Debug(lg)
		}
	}
	l.FlushBuf()
	return nil
}

func NewLogWriteCloser(bid int32) *LogWriteCloser {
	return &LogWriteCloser{
		bid: bid,
	}
}
