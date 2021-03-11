package avro

import (
	"bytes"

	"github.com/linkedin/goavro"
)

type Avro struct {
	ocfr   *goavro.OCFReader
	codec  *goavro.Codec
	schema string
	data   map[string]interface{}
}

func NewAvroReader(data []byte) (*Avro, error) {
	o, err := goavro.NewOCFReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	readCodec, err := goavro.NewCodec(o.Codec().Schema())
	if err != nil {
		return nil, err
	}
	return &Avro{
		ocfr:   o,
		codec:  readCodec,
		schema: o.Codec().Schema(),
	}, nil
}

func (a *Avro) AvroToJson() (string, error) {
	bytes, err := a.AvroToByteString()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (a *Avro) AvroToByteString() ([]byte, error) {
	if a.data == nil {
		var err error
		a.data, err = a.AvroToMap()
		if err != nil {
			return nil, err
		}
	}
	jbytes, err := a.codec.TextualFromNative(nil, a.data)
	if err != nil {
		return nil, err
	}
	return jbytes, nil
}

func (a *Avro) AvroToMap() (map[string]interface{}, error) {
	if a.data == nil {
		for a.ocfr.Scan() {
			datum, err := a.ocfr.Read()
			if err != nil {
				return nil, err
			}
			a.data = datum.(map[string]interface{})
		}
	}
	return a.data, nil
}
