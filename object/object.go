package object

import "fmt"

// / ObjectType enumerates the types of objects that our evaluator will work
// / with.
type ObjectType uint

const (
	O_INTEGER = iota
	O_BOOLEAN
	O_NULL
)

// String returns a mostly-human-readable string enum value for the
// corresponding ObjectType.
func (ot ObjectType) String() string {
	switch ot {
	case O_INTEGER:
		return "INTEGER"
	case O_BOOLEAN:
		return "BOOLEAN"
	case O_NULL:
		return "NULL"
	default:
		return "UNKNOWN"
	}
}

// Object is the return value type of all Monkey expressions.
type Object interface {
	// Type returns this object's type.
	Type() ObjectType
	// Inspect gives an repl-friendly output rendering of this object.
	Inspect() string
}

// Integer is an object that represents a 64-bit signed integer.
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return O_INTEGER
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Integer is an object that represents a 64-bit signed integer.
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return O_BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null is an object that represents ... null.
type Null struct{}

func (n *Null) Type() ObjectType {
	return O_NULL
}

func (n *Null) Inspect() string {
	return "null"
}
