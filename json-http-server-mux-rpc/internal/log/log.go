package log

import (
	"encoding/json"
	"net/http"
	"sync"
)

// Log store
type Log struct {
	Mu      sync.Mutex
	Entries []Record
}

// Single (value + offset)
type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"`
}

type AppendRequest struct {
	Record Record `json:"record"`
}

type AppendResponse struct {
	Offset uint64 `json:"offset"`
}

type FetchRequest struct {
	Offset uint64 `json:"offset"`
}

type FetchResponse struct {
	Record Record `json:"record"`
}

// Add to log
func (l *Log) Append(w http.ResponseWriter, r *http.Request) {
	var req AppendRequest
	// Decode JSON body into AppendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.Mu.Lock()
	defer l.Mu.Unlock()

	// Assign offset and append entry
	req.Record.Offset = uint64(len(l.Entries))
	l.Entries = append(l.Entries, req.Record)

	// Encode and return response with new offset
	res := AppendResponse{Offset: req.Record.Offset}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Retrieve from log
func (l *Log) Fetch(w http.ResponseWriter, r *http.Request) {
	var req FetchRequest
	// Decode JSON body into FetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l.Mu.Lock()
	defer l.Mu.Unlock()

	// Validate offset
	if req.Offset >= uint64(len(l.Entries)) {
		http.Error(w, "offset out of range", http.StatusNotFound)
		return
	}

	// Encode and return the requested record
	res := FetchResponse{Record: l.Entries[req.Offset]}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
