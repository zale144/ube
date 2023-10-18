package model

import "encoding/json"

// Request is the generic API request
type Request struct {
	Path       string `json:"path,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	RequestID  string `json:"request_id,omitempty"`
	Body       string `json:"body,omitempty"`
	SourceURI  string `json:"source_uri"`
	Reference
}

// Reference holds headers, query and path params
type Reference struct {
	Headers               map[string]string `json:"headers,omitempty"`
	QueryStringParameters map[string]string `json:"query_string_parameters,omitempty"`
	PathParameters        map[string]string `json:"path_parameters,omitempty"`
}

// GetID returns the request ID (implements model.Input)
func (r Request) GetID() string {
	return r.RequestID
}

// GetReference returns the request reference (implements model.Input)
func (r Request) GetReference() string {
	jsnRef, _ := json.Marshal(r.Reference)
	return string(jsnRef)
}

// GetBody returns the request body (implements model.Input)
func (r Request) GetBody() string {
	return r.Body
}

// GetSourceURI returns the request source URI (implements model.Input)
func (r Request) GetSourceURI() string {
	return r.SourceURI
}

// Response is the generic API response
type Response struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}
