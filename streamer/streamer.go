package streamer

import (
	"bufio"
	"fmt"
	"io"
)

type LogStreamer struct {
	reader io.Reader
}

func NewLogStreamer(r io.Reader) (*LogStreamer, error) {
	if r == nil {
		return nil, fmt.Errorf("reader is invalid")
	}
	return &LogStreamer{reader: r}, nil
}

func (lg *LogStreamer) EachLine(fn func(line string) error) error {

	scanner := bufio.NewScanner(lg.reader)

	maxBufferSize := 1024 * 1024 // 1 MB
	buf := make([]byte, maxBufferSize)
	scanner.Buffer(buf, maxBufferSize)

	for scanner.Scan() {
		line := scanner.Text()

		if err := fn(line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}
