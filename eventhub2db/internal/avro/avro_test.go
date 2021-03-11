package avro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

func prepare() ([]byte, error) {
	// avro schema defintion
	codec, err := goavro.NewCodec(`
	{
		"type": "record",
		"name": "my.test.name",
		"doc": "test.doc",
		"fields" : [
		{
			"name": "test_string",
			"type": "string"
		},
		{
			"name": "test_int",
			"type": "int"
		},
		{
			"name": "test_double",
			"type": "double"
		},
		{
			"name": "test_bool",
			"type": "boolean"
		},
		{
			"name": "test_long",
			"type": "long"
		}
		]
	}`)
	if err != nil {
		fmt.Println(err)
	}

	msg := make(map[string]interface{})

	msg["test_string"] = "myTestString"
	msg["test_int"] = 123456
	msg["test_bool"] = true
	msg["test_double"] = 3.145678
	msg["test_long"] = 654654

	bin := new(bytes.Buffer)
	if err != nil {
		return nil, err
	}

	ocfw, err := goavro.NewOCFWriter(goavro.OCFConfig{
		W:     bin,
		Codec: codec,
	})
	if err != nil {
		return nil, err
	}

	err = ocfw.Append([]interface{}{msg})
	if err != nil {
		return nil, err
	}
	return bin.Bytes(), nil
}

func TestAvroToByteString(t *testing.T) {
	assert := assert.New(t)
	bytes, err := prepare()
	assert.Nil(err)
	o, err := NewAvroReader(bytes)
	assert.Nil(err)
	m, err := o.AvroToByteString()
	assert.Nil(err)
	var dat map[string]interface{}
	d := json.NewDecoder(strings.NewReader(string(m)))
	d.UseNumber()
	if err := d.Decode(&dat); err != nil {
		log.Fatal(err)
	}

	// Raw json decoding uses json.Number vor numbers. json.Number cannot be int32.
	assert.Equal(dat["test_string"].(string), "myTestString")
	assert.Equal(dat["test_bool"].(bool), true)
	double, err := dat["test_double"].(json.Number).Float64()
	assert.Nil(err)
	assert.Equal(double, float64(3.145678))
	timestamp, err := dat["test_int"].(json.Number).Int64()
	assert.Nil(err)
	assert.Equal(timestamp, int64(123456))
	long, err := dat["test_long"].(json.Number).Int64()
	assert.Nil(err)
	assert.Equal(long, int64(654654))
}

func TestAvroToByteStringTwice(t *testing.T) {
	assert := assert.New(t)
	bytes, err := prepare()
	assert.Nil(err)
	o, err := NewAvroReader(bytes)
	assert.Nil(err)
	m1, err := o.AvroToByteString()
	assert.Nil(err)
	var dat1 map[string]interface{}
	d1 := json.NewDecoder(strings.NewReader(string(m1)))
	d1.UseNumber()
	if err := d1.Decode(&dat1); err != nil {
		log.Fatal(err)
	}

	// Raw json decoding uses json.Number vor numbers. json.Number cannot be int32.
	assert.Equal(dat1["test_string"].(string), "myTestString")
	assert.Equal(dat1["test_bool"].(bool), true)
	double, err := dat1["test_double"].(json.Number).Float64()
	assert.Nil(err)
	assert.Equal(double, float64(3.145678))
	timestamp, err := dat1["test_int"].(json.Number).Int64()
	assert.Nil(err)
	assert.Equal(timestamp, int64(123456))
	long, err := dat1["test_long"].(json.Number).Int64()
	assert.Nil(err)
	assert.Equal(long, int64(654654))

	// Second time, because the internal functions within the used
	// goavro library is intended to read only once from the OCF container.
	// This must be working to show the buffering
	m2, err := o.AvroToByteString()
	assert.Nil(err)
	var dat2 map[string]interface{}
	d2 := json.NewDecoder(strings.NewReader(string(m2)))
	d2.UseNumber()
	if err := d2.Decode(&dat2); err != nil {
		log.Fatal(err)
	}
	assert.Equal(dat2["test_string"].(string), "myTestString")
	assert.Equal(dat2["test_bool"].(bool), true)
	double2, err := dat2["test_double"].(json.Number).Float64()
	assert.Nil(err)
	assert.Equal(double2, float64(3.145678))
	timestamp2, err := dat2["test_int"].(json.Number).Int64()
	assert.Nil(err)
	assert.Equal(timestamp2, int64(123456))
	long2, err := dat2["test_long"].(json.Number).Int64()
	assert.Nil(err)
	assert.Equal(long2, int64(654654))

	assert.Equal(dat1, dat2)
}

func TestAvroToJson(t *testing.T) {
	assert := assert.New(t)
	bytes, err := prepare()
	assert.Nil(err)
	o, err := NewAvroReader(bytes)
	assert.Nil(err)
	_, err = o.AvroToJson()
	assert.Nil(err)
}

func TestAvroToMapTwice(t *testing.T) {
	assert := assert.New(t)
	bytes, err := prepare()
	assert.Nil(err)
	o, err := NewAvroReader(bytes)
	assert.Nil(err)
	a, err := o.AvroToMap()
	assert.Nil(err)
	b, err := o.AvroToMap()
	assert.Nil(err)
	assert.Equal(a, b)
}

func TestAvroToMap(t *testing.T) {
	assert := assert.New(t)
	bytes, err := prepare()
	assert.Nil(err)
	o, err := NewAvroReader(bytes)
	assert.Nil(err)
	m, err := o.AvroToMap()
	assert.Nil(err)
	assert.Equal(m["test_string"].(string), "myTestString")
	assert.Equal(m["test_int"].(int32), int32(123456))
	assert.Equal(m["test_bool"].(bool), true)
	assert.Equal(m["test_long"].(int64), int64(654654))
	assert.Equal(m["test_double"].(float64), 3.145678)
}
