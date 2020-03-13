package model

// RequestType describes a single request type.
type RequestType struct {
	// the request type identifier
	ID string `json:"id"`

	// the request type name
	Name string `json:"name"`
}

// RequestTypeListing describes a listing of request types.
type RequestTypeListing struct {
	// the request types
	RequestTypes []*RequestType `json:"request_types"`
}
