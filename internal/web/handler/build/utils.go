package build

type LogWriter struct {
	bid int32
}

func (l *LogWriter) Write(p []byte) (n int, err error) {
	createLog(l.bid, string(p))
	return len(p), nil
}

func NewLogWriteCloser(bid int32) *LogWriter {
	return &LogWriter{
		bid: bid,
	}
}
