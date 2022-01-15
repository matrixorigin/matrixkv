package metadata

const (
	// SetType set operation type
	SetType uint64 = 1
	// GetType get operation type
	GetType uint64 = 2
	// DeleteType delete operation type
	DeleteType uint64 = 3
)

// SetRequest set request
type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SetResponse set response
type SetResponse struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

// GetRequest get request
type GetRequest struct {
	Key string `json:"key"`
}

// GetResponse set response
type GetResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Error string `json:"error"`
}

// DeleteRequest delete request
type DeleteRequest struct {
	Key string `json:"key"`
}

// DeleteResponse delete response
type DeleteResponse struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}
