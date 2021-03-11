package avro

import (
	"bytes"

	"github.com/linkedin/goavro"
)

type Avro struct {
	ocfr   *goavro.OCFReader
	codec  *goavro.Codec
	schema string
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
	m, err := a.AvroToMap()
	if err != nil {
		return nil, err
	}
	jbytes, err := a.codec.TextualFromNative(nil, m)
	if err != nil {
		return nil, err
	}
	return jbytes, nil
}

func (a *Avro) AvroToMap() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for a.ocfr.Scan() {
		datum, err := a.ocfr.Read()
		if err != nil {
			return nil, err
		}
		m = datum.(map[string]interface{})
	}
	return m, nil
}
