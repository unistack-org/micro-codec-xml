package xml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "with replacement",
			input:    []byte(`<?xml version="1.1" encoding="UTF-8"?><root name="NAME"><item>ITEM</item></root>`),
			expected: []byte(`<?xml version="1.0" encoding="UTF-8"?><root name="NAME"><item>ITEM</item></root>`),
		},
		{
			name:     "without replacement",
			input:    []byte(`<?xml version="1.0" encoding="UTF-8"?><root name="NAME"><item>ITEM</item></root>`),
			expected: []byte(`<?xml version="1.0" encoding="UTF-8"?><root name="NAME"><item>ITEM</item></root>`),
		},
		{
			name:     "check invalid replacement",
			input:    []byte(`<?xml version="1.0" encoding="UTF-8"?><root version="1.1"><item>ITEM</item></root>`),
			expected: []byte(`<?xml version="1.0" encoding="UTF-8"?><root version="1.1"><item>ITEM</item></root>`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, len(tt.input))
			n, err := newReader(tt.input).Read(buf)
			require.NoError(t, err)
			require.Equal(t, tt.expected, buf[:n])
		})
	}
}
