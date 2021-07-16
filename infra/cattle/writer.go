package cattle

import "github.com/go-xorm/builder"

type BytesWriter struct {
	writer *builder.StringBuilder
	args   []interface{}
}

// NewWriter creates a new string writer
func NewWriter() *BytesWriter {
	w := &BytesWriter{
		writer: &builder.StringBuilder{},
	}
	return w
}

// Write writes data to Writer
func (s *BytesWriter) Write(buf []byte) (int, error) {
	return s.writer.Write(buf)
}

// Append appends args to Writer
func (s *BytesWriter) Append(args ...interface{}) {
	s.args = append(s.args, args...)
}
