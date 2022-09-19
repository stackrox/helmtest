package parser

import "fmt"

// SourceContext
type SourceContext struct {
	Filename string
	Line     int
}

func (c SourceContext) String() string {
	filename := c.Filename
	if filename == "" {
		filename = "<input>"
	}
	return fmt.Sprintf("%s:%d", filename, c.Line+1)
}

func (c *SourceContext) IsZero() bool {
	return c.Filename != "" || c.Line != 0
}
