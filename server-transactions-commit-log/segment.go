package log

import (
	"fmt"
	"math"
	"os"
	"path"

	log_v1 "server-transactions-commit-log/api/v1"

	"google.golang.org/protobuf/proto"
)

type segment struct {
	store                  *store
	index                  *index
	baseOffset, nextOffset uint64
	config                 Config
}

func newSegment(dir string, baseOffset uint64, c Config) (*segment, error) {
	s := &segment{
		baseOffset: baseOffset,
		config:     c,
	}

	var err error

	// Create the store file
	storeFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".store")),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.store, err = newStore(storeFile); err != nil {
		return nil, err
	}

	// Create the index file
	indexFile, err := os.OpenFile(
		path.Join(dir, fmt.Sprintf("%d%s", baseOffset, ".index")),
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return nil, err
	}
	if s.index, err = newIndex(indexFile, c); err != nil {
		return nil, err
	}

	// Calculate the next offset
	if off, _, err := s.index.Read(int64(math.MaxInt64)); err != nil { // Conversion to int64
		s.nextOffset = baseOffset
	} else {
		s.nextOffset = baseOffset + uint64(off) + 1 // Convert off to uint64
	}

	return s, nil
}

func (s *segment) Append(record *log_v1.Record) (offset uint64, err error) {
	// Get the current offset
	cur := s.nextOffset
	record.Offset = cur

	// Serialize the record using protobuf
	p, err := proto.Marshal(record)
	if err != nil {
		return 0, err
	}

	// Write the serialized record to the store
	_, pos, err := s.store.Append(p)
	if err != nil {
		return 0, err
	}

	// Write to the index
	if err = s.index.Write(
		// Index offsets are relative to the base offset on the store file
		uint32(s.nextOffset-s.baseOffset), // Convert uint64 to uint32
		uint64(pos),                       // Convert pos to uint64
	); err != nil {
		return 0, err
	}

	// Increment the next available offset
	s.nextOffset++

	return cur, nil
}

func (s *segment) Read(off uint64) (*log_v1.Record, error) {
	// Transform the absolute offset into a relative offset
	_, pos, err := s.index.Read(int64(off - s.baseOffset)) // Convert to int64
	if err != nil {
		return nil, err
	}

	// Read the record from the store using the obtained position
	p, err := s.store.Read(uint64(pos)) // Convert pos to uint64
	if err != nil {
		return nil, err
	}

	// Deserialize the record from protobuf
	var record log_v1.Record
	if err := proto.Unmarshal(p, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (s *segment) IsMaxed() bool {
	// Check if the store or index file has reached the maximum size allowed
	if s.store.size >= s.config.Segment.MaxStoreBytes || s.index.size >= s.config.Segment.MaxIndexBytes {
		return true
	}
	return false
}

func (s *segment) Remove() error {
	// Close the files before attempting to remove them
	if err := s.index.Close(); err != nil {
		return err
	}
	if err := s.store.Close(); err != nil {
		return err
	}

	// Remove the index file
	if err := os.Remove(s.index.Name()); err != nil {
		return err
	}

	// Remove the store file
	if err := os.Remove(s.store.Name()); err != nil {
		return err
	}

	return nil
}

func (s *segment) Close() error {
	var err error

	// Close the index file
	if cerr := s.index.Close(); cerr != nil {
		err = cerr
	}

	// Close the store file
	if serr := s.store.Close(); serr != nil {
		if err != nil {
			err = serr
		}
	}

	return err
}
