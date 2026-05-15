package xml

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
	"go.unistack.org/micro/v5/codec"
)

func TestXmlCodec_Unmarshal(t *testing.T) {
	type result struct {
		XMLName xml.Name `xml:"root"`
		Name    string   `xml:"name,attr"`
		Item    string   `xml:"item"`
	}

	tests := []struct {
		name     string
		input    []byte
		isXML11  bool
		expected result
	}{
		{
			name:  "without version",
			input: []byte(`<root name="NAME"><item>ITEM</item></root>`),
			expected: result{
				XMLName: xml.Name{Local: "root"},
				Name:    "NAME",
				Item:    "ITEM",
			},
		},
		{
			name:  "version 1.0",
			input: []byte(`<?xml version="1.0" encoding="UTF-8"?><root name="NAME"><item>ITEM</item></root>`),
			expected: result{
				XMLName: xml.Name{Local: "root"},
				Name:    "NAME",
				Item:    "ITEM",
			},
		},
		{
			name:    "version 1.1",
			input:   []byte(`<?xml version="1.1" encoding="UTF-8"?><root name="NAME"><item>ITEM</item></root>`),
			isXML11: true,
			expected: result{
				XMLName: xml.Name{Local: "root"},
				Name:    "NAME",
				Item:    "ITEM",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				r    result
				opts []codec.Option
			)

			if tt.isXML11 {
				opts = append(opts, UnmarshalXML11(true))
			}

			err := NewCodec(opts...).Unmarshal(tt.input, &r)
			require.NoError(t, err)
			require.Equal(t, tt.expected, r)
		})
	}
}
