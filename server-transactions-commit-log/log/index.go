package log

import (
	"encoding/binary"
	"io"
	"os"
	"sync"

	"github.com/tysonmote/gommap"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth        = int(offWidth + posWidth)
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
	mu   sync.Mutex
}

// newIndex creates a new index from a file
func newIndex(f *os.File, c Config) (*index, error) {
	// 1. Obtain the size of the file we are indexing
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	idx := &index{
		file: f,
		size: uint64(fi.Size()),
	}

	// 2. Set the size of our index as the size of the file
	if err := os.Truncate(f.Name(), int64(c.Segment.MaxIndexBytes)); err != nil {
		return nil, err
	}

	// 3. Map the file directly to memory
	if idx.mmap, err = gommap.Map(idx.file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED); err != nil {
		return nil, err
	}

	// Zeroing the mmap if the size is 0 to ensure no garbage data
	if len(idx.mmap) > 0 && idx.size == 0 {
		for i := range idx.mmap {
			idx.mmap[i] = 0
		}
	}

	return idx, nil
}

// Write writes an offset and position to the index
func (i *index) Write(off uint32, pos uint64) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 1. First, get the size of our index and validate that we have space to write
	if i.size+uint64(entWidth) > uint64(len(i.mmap)) {
		return io.EOF
	}

	// 2. Write the offset first, from the end of the file up to the offset size
	binary.BigEndian.PutUint32(i.mmap[i.size:], off)

	// 3. Then write the position from the end of the offset to the end of the position
	binary.BigEndian.PutUint64(i.mmap[i.size+uint64(offWidth):], pos)

	i.size += uint64(entWidth)

	return nil
}

// Read reads an entry from the index
func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 1. Simply get the place where we want to read
	if i.size == 0 {
		return 0, 0, io.EOF
	}

	if in == -1 {
		out = uint32((i.size / uint64(entWidth)) - 1)
		offsetPos := i.size - uint64(entWidth)

		// 2. Multiply the integer by entWidth, giving us a starting position to decode from binary
		pos = binary.BigEndian.Uint64(i.mmap[offsetPos+uint64(offWidth):])

		return out, pos, nil
	}

	entryOffset := uint64(in) * uint64(entWidth)
	if in < 0 || entryOffset >= i.size {
		return 0, 0, io.EOF
	}

	// 2.1 Start from the position we established earlier, then get from there to the end of the offset bytes
	out = binary.BigEndian.Uint32(i.mmap[entryOffset:])

	// 2.2 For the position, simply start from the end of the offset to the entWidth, this will give us the Store position
	pos = binary.BigEndian.Uint64(i.mmap[entryOffset+uint64(offWidth):])

	return out, pos, nil
}

// Close closes the index file
func (i *index) Close() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 1. First, check that we have flushed the data from memory to the file
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}

	// Force file sync to ensure data is written to disk
	if err := i.file.Sync(); err != nil {
		return err
	}

	// 2. Then truncate the file to the exact size of its content
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}

	// 3. Finally, simply close the file
	return i.file.Close()
}

// Name returns the name of the index file
func (i *index) Name() string {
	return i.file.Name()
}
