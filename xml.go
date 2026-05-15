// Package xml provides a xml codec
package xml

import (
	"bytes"
	"encoding/xml"
	"io"

	pb "go.unistack.org/micro-proto/v5/codec"
	"go.unistack.org/micro/v5/codec"
	rutil "go.unistack.org/micro/v5/util/reflect"
)

var _ codec.Codec = &xmlCodec{}

type xmlCodec struct {
	opts codec.Options
}

const (
	flattenTag = "flatten"
)

type unmarshalStrictKey struct{}

func UnmarshalStrict(b bool) codec.Option {
	return codec.SetOption(unmarshalStrictKey{}, b)
}

type unmarshalXML11Key struct{}

func UnmarshalXML11(b bool) codec.Option {
	return codec.SetOption(unmarshalXML11Key{}, b)
}

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

	var unmarshalStrict bool
	var unmarshalXML11 bool
	if options.Context != nil {
		if v, ok := options.Context.Value(unmarshalStrictKey{}).(bool); ok {
			unmarshalStrict = v
		}
		if v, ok := options.Context.Value(unmarshalXML11Key{}).(bool); ok {
			unmarshalXML11 = v
		}
	}

	var r io.Reader

	if unmarshalXML11 {
		r = newReader(b)
	} else {
		r = bytes.NewReader(b)
	}

	d := xml.NewDecoder(r)

	if unmarshalStrict {
		d.Strict = true
	}

	return d.Decode(v)
}

type reader struct {
	s        []byte
	i        int64 // current reading index
	prevRune int   // index of previous rune; or < 0
	skip     bool
}

func newReader(b []byte) io.Reader {
	return &reader{b, 0, -1, false}
}

var (
	srcP = []byte(`<?xml version="1.1"`)
	srcL = len(srcP)
	dstP = []byte(`<?xml version="1.0"`)
)

func (r *reader) Read(b []byte) (n int, err error) {
	if r.i >= int64(len(r.s)) {
		return 0, io.EOF
	}
	r.prevRune = -1
	n = copy(b, r.s[r.i:])
	if !r.skip {
		if idx := bytes.Index(b, srcP); idx >= 0 {
			copy(b[idx:idx+srcL], dstP)
			r.skip = true
		}
	}

	r.i += int64(n)
	return
}

func (c *xmlCodec) String() string {
	return "xml"
}

func NewCodec(opts ...codec.Option) *xmlCodec {
	return &xmlCodec{opts: codec.NewOptions(opts...)}
}
