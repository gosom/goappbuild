package goappbuild

import (
	"time"
)

// AttributeType is a string that represents the type of an attribute
type AttributeType string

const (
	// AttributeTypeString is the string type
	AttributeTypeString AttributeType = "string"
	// AttributeTypeNumber is the number type
	AttributeTypeInteger AttributeType = "integer"
	// AttributeTypeNumber is the number type
	AttributeTypeNumeric AttributeType = "numeric"
	// AttributeTypeBoolean is the boolean type
	AttributeTypeFloat AttributeType = "float"
	// AttributeTypeBoolean is the boolean type
	AttributeTypeBoolean AttributeType = "boolean"
	// AttributeTypeDate is the date type
	AttributeTypeTime AttributeType = "time"
	// AttributeTypeUUID is the UUID type
	AttributeTypeUUID AttributeType = "uuid"
	// AttributeTypeJSON is the JSON type
	AttributeTypeJSON AttributeType = "json"
)

// Attribute is a struct that represents an attribute of a collection
type Attribute struct {
	// Name is the name of the attribute
	Name string
	// Type is the type of the attribute
	Type AttributeType
	// Required is a boolean that indicates if the attribute is required
	Required bool
	// Unique is a boolean that indicates if the attribute is unique
	Unique bool
	// Primary is a boolean that indicates if the attribute is the primary key
	Primary bool
	// Index is a boolean that indicates if the attribute is indexed
	Index bool
	// Relationships holds the relationships of the attribute
	Relationships []Relationship
	// CreatedAt is the time the attribute was created
	CreatedAt time.Time
	// UpdatedAt is the time the attribute was last updated
	UpdatedAt time.Time
}

// Relationship is a struct that represents a relationship of an attribute
type Relationship struct {
	Reference string
}
