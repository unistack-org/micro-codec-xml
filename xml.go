// Package xml provides a xml codec
package xml

import (
	"encoding/xml"
	"io"
	"io/ioutil"

	"github.com/unistack-org/micro/v3/codec"
)

type xmlCodec struct{}

func (c *xmlCodec) Marshal(b interface{}) ([]byte, error) {
	switch m := b.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	}

	return xml.Marshal(b)
}

func (c *xmlCodec) Unmarshal(b []byte, v interface{}) error {
	if b == nil {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = b
		return nil
	}

	return xml.Unmarshal(b, v)
}

func (c *xmlCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *xmlCodec) ReadBody(conn io.Reader, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		m.Data = buf
		return nil
	}

	err := xml.NewDecoder(conn).Decode(b)
	if err == io.EOF {
		return nil
	}
	return err
}

func (c *xmlCodec) Write(conn io.Writer, m *codec.Message, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	}

	return xml.NewEncoder(conn).Encode(b)
}

func (c *xmlCodec) String() string {
	return "xml"
}

func NewCodec() codec.Codec {
	return &xmlCodec{}
}
