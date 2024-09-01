package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

type store struct {
	File *os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())
	return &store{
		File: f,
		buf:  bufio.NewWriter(f),
		size: size,
	}, nil
}

// Append adds data to the store
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Receive the data to be written, in this case, bytes
	pos = s.size

	// 2. Determine the position where the data will be written
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}

	// 3. First, write to a buffer, then write to the physical file
	if err := binary.Write(s.buf, enc, p); err != nil {
		return 0, 0, err
	}

	// 4. After writing, obtain the new length of the file
	n = uint64(len(p)) + uint64(lenWidth)
	s.size += n

	// 5. Return the number of bytes written and the current position of the file
	return n, pos, nil
}

// Read retrieves data from the store at the given position
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. As we are dealing with files, flush any data that hasn't been written to the file yet and is still in the buffer
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	// 2. Determine how many bytes we need to read to reach our store
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	// 3. Finally, retrieve the store from the desired position by obtaining the number of bytes needed to traverse the file
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}
	return b, nil
}

// ReadAt is a helper function to make it easier to retrieve the necessary bytes to read a store
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure the buffer is flushed before reading
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

// Close closes the store, persisting any remaining data in the buffer to the file
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure all data is written to the file before closing
	if err := s.buf.Flush(); err != nil {
		return err
	}
	return s.File.Close()
}

func (s *store) Name() string {
    return s.File.Name()
}
