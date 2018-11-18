package encodedbytes

import (
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"bytes"
	"io/ioutil"
)

type (
	Converter struct {
		from string
		to string
	}
)

func resolveEncoding(b []byte, name string) encoding.Encoding {
	switch name {
	case "UTF-8":
		return nil
	case "UTF-16":
		if b == nil {
			return unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
		} else if len(b) < 2 {
			return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		} else if b[0] == 0xFE && b[1] == 0xFF {
			return unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
		} else if b[0] == 0xFF && b[1] == 0xFE {
			return unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
		} else {
			return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		}
	case "UTF-16LE":
		return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	case "UTF-16BE":
		return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	case "ISO-8859-1":
		return charmap.Windows1252
	}
	return nil
}

func (c Converter) ConvertString(s string) (string, error) {
	var native []byte

	b := []byte(s)

	fromEncoding := resolveEncoding(b, c.from)
	if fromEncoding != nil {
		in := bytes.NewReader(b)
		reader := transform.NewReader(in, fromEncoding.NewDecoder())
		n, err := ioutil.ReadAll(reader)
		if err != nil {
			return "", err
		}
		native = n
	} else {
		native = b
	}

	toEncoding := resolveEncoding(nil, c.to)
	if toEncoding != nil {
		out := &bytes.Buffer{}
		writer := transform.NewWriter(out, toEncoding.NewEncoder())
		defer writer.Close()
		_, err := writer.Write(native)
		if err != nil {
			return "", err
		}
		return string(out.Bytes()), nil
	} else {
		return string(native), nil
	}
}

func NewConverter(from string, to string) (*Converter, error) {
	return &Converter{from:from, to:to}, nil
}
