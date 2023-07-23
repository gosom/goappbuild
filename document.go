package goappbuild

import "encoding/json"

// Document is a struct that represents a document
type Document struct {
	Values map[string]any
}

func (d Document) MarshalJSON() ([]byte, error) {
	if d.Values == nil {
		d.Values = make(map[string]any)
	}

	return json.Marshal(d.Values)
}
