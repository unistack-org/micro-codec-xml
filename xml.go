// Package xml provides a xml codec
package xml // import "go.unistack.org/micro-codec-xml/v3"

import (
	"encoding/xml"

	pb "go.unistack.org/micro-proto/v3/codec"
	"go.unistack.org/micro/v3/codec"
	rutil "go.unistack.org/micro/v3/util/reflect"
)

var _ codec.Codec = &xmlCodec{}

type xmlCodec struct {
	opts codec.Options
}

const (
	flattenTag = "flatten"
)

func (c *xmlCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
	}

	return xml.Marshal(v)
}

func (c *xmlCodec) Unmarshal(b []byte, v interface{}, opts ...codec.Option) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = b
		return nil
	case *pb.Frame:
		m.Data = b
		return nil
	case *codec.RawMessage:
		*m = append((*m)[0:0], b...)
		return nil
	case codec.RawMessage:
		copy(m, b)
		return nil
	}

	return xml.Unmarshal(b, v)
}

func (c *xmlCodec) String() string {
	return "xml"
}

func NewCodec(opts ...codec.Option) *xmlCodec {
	return &xmlCodec{opts: codec.NewOptions(opts...)}
}
