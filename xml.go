// Package xml provides a xml codec
package xml // import "go.unistack.org/micro-codec-xml/v3"

import (
	"encoding/xml"
	"io"

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
	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
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

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = b
		return nil
	case *pb.Frame:
		m.Data = b
		return nil
	}

	return xml.Unmarshal(b, v)
}

func (c *xmlCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *xmlCodec) ReadBody(conn io.Reader, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}

	return c.Unmarshal(buf, v)
}

func (c *xmlCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := c.Marshal(v)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return codec.ErrInvalidMessage
	}

	_, err = conn.Write(buf)
	return err
}

func (c *xmlCodec) String() string {
	return "xml"
}

func NewCodec(opts ...codec.Option) *xmlCodec {
	return &xmlCodec{opts: codec.NewOptions(opts...)}
}
