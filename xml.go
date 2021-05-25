// Package xml provides a xml codec
package xml

import (
	"encoding/xml"
	"io"

	"github.com/unistack-org/micro/v3/codec"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

type xmlCodec struct{}

const (
	flattenTag = "flatten"
)

func (c *xmlCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	return xml.Marshal(v)
}

func (c *xmlCodec) Unmarshal(b []byte, v interface{}) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	if m, ok := v.(*codec.Frame); ok {
		m.Data = b
		return nil
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
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

func NewCodec() codec.Codec {
	return &xmlCodec{}
}
