package models

// Employee represents an employee entity matching the MySQL table schema.
//
// GOLANG CONCEPT - STRUCTS:
// In Golang, a "struct" (structure) is a custom data type that groups together
// related variables (fields) under one name. It is similar to a class or object in
// other languages like Java, C#, or JavaScript, but without inheritance.
//
// GOLANG CONCEPT - JSON STRUCT TAGS:
// The backtick annotations like `json:"id"` are struct tags.
// Struct tags provide metadata about fields. The `json:"..."` tag instructs
// Go's encoding/json package (and Gin's ShouldBindJSON) how to map Go struct fields
// to JSON keys when encoding (sending JSON) or decoding (parsing incoming JSON).
// - "id" in JSON maps to "ID" in Go
// - "name" in JSON maps to "Name" in Go
// - "npwp" in JSON maps to "NPWP" in Go
// - "address" in JSON maps to "Address" in Go
type Employee struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	NPWP    string `json:"npwp"`
	Address string `json:"address"`
}
