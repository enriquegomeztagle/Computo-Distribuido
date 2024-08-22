package log

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Log -> many records
type Log struct {
	mu      sync.Mutex
	entries []Record
}

// Record -> 1 log (value + offset)
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

// AppendRequest -> receives JSON -> log entry
type AppendRequest struct {
	Record Record `json:"record"`
}

// AppendResponse -> offset
type AppendResponse struct {
	Offset uint64 `json:"offset"`
}

// FetchRequest -> receives JSON -> consumes with offset
type FetchRequest struct {
	Offset uint64 `json:"offset"`
}

// FetchResponse -> response JSON
type FetchResponse struct {
	Record Record `json:"record"`
}

// Append -> add record to log
func (l *Log) Append(w http.ResponseWriter, r *http.Request) {
	var req AppendRequest
	// Decode JSON body into AppendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Assign offset and append entry
	req.Record.Offset = uint64(len(l.entries))
	l.entries = append(l.entries, req.Record)

	// Encode and return response with new offset
	res := AppendResponse{Offset: req.Record.Offset}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Fetch -> retrieves record according to offset
func (l *Log) Fetch(w http.ResponseWriter, r *http.Request) {
	var req FetchRequest
	// Decode JSON body into FetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Validate offset
	if req.Offset >= uint64(len(l.entries)) {
		http.Error(w, "offset out of range", http.StatusNotFound)
		return
	}

	// Encode and return response with requested record
	res := FetchResponse{Record: l.entries[req.Offset]}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
