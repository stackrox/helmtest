package parser

import "fmt"

// SourceContext stores information about location in a source file.
type SourceContext struct {
	Filename string
	Line     int
}

// String returns a human-readable string representation of this source context.
func (c SourceContext) String() string {
	filename := c.Filename
	if filename == "" {
		filename = "<input>"
	}
	return fmt.Sprintf("%s:%d", filename, c.Line+1) // c.Line is zero-based, add 1 for human-readable
}

// IsZero checks whether this SourceContext is equal to the zero SourceContext.
func (c *SourceContext) IsZero() bool {
	return c.Filename == "" && c.Line == 0
}
