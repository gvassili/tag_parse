package tag_parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertTagParam(t *testing.T, tag string, expectedParams []Param) {
	params, err := Parse(tag)
	assert.NoError(t, err, fmt.Sprintf("tag=`%s`", tag))
	for i, param := range params {
		assert.Equal(t, expectedParams[i].Key, param.Key, fmt.Sprintf("tag=`%s`", tag))
		assert.Equal(t, expectedParams[i].Values, param.Values, fmt.Sprintf("tag=`%s`", tag))
	}
}

func assertTagParamError(t *testing.T, tag string) {
	_, err := Parse(tag)
	assert.Error(t, err, fmt.Sprintf("tag=`%s`", tag))
}


func TestParse(t *testing.T) {
	assertTagParam(t, ``, nil)
	assertTagParam(t, `key1`, []Param{{"key1", nil}})
	assertTagParam(t, `key1=`, []Param{{"key1", nil}})
	assertTagParam(t, `key1,key2`, []Param{{"key1", nil}, {"key2", nil}})
	assertTagParam(t, `key1=value1`, []Param{{"key1", []string{"value1"}}})
	assertTagParam(t, `' key1 '=' value1 '`, []Param{{" key1 ", []string{" value1 "}}})
	assertTagParam(t, `key1='value1,value2'`, []Param{{"key1", []string{"value1,value2"}}})
	assertTagParam(t, `key1=value1,key2=value2`, []Param{{"key1", []string{"value1"}}, {"key2", []string{"value2"}}})
	assertTagParam(t, `key1 = value1 , key2 = value2`, []Param{{"key1", []string{"value1"}}, {"key2", []string{"value2"}}})
	assertTagParam(t, `key1 = value1 value2 value3`, []Param{{"key1", []string{"value1", "value2", "value3"}}})
	assertTagParam(t, `key1 = 'value 1' 'value 2' 'value 3'`, []Param{{"key1", []string{"value 1", "value 2", "value 3"}}})
	assertTagParam(t, `key1 = 'value 1''value 1''value 1'`, []Param{{"key1", []string{"value 1value 1value 1"}}})
	assertTagParam(t, `key1='\\\'value \'1' ' \\value\' 2'`, []Param{{"key1", []string{`\'value '1`, ` \value' 2`}}})
	assertTagParam(t, `key1=\\\'value\'1 \\value\'2`, []Param{{"key1", []string{`\'value'1`, `\value'2`}}})

	assertTagParamError(t, `key1 key2`)
	assertTagParamError(t, `key1 key2 = value1`)
	assertTagParamError(t, `'key1 ' 'key 2'`)
	assertTagParamError(t, `key1\n`)
	assertTagParamError(t, `key1=value1\n`)
	assertTagParamError(t, `'key1`)
	assertTagParamError(t, `key1='value1`)
	assertTagParamError(t, `'key1=value1`)
	assertTagParamError(t, `'key1='value1'`)
	assertTagParamError(t, `key1=value1\`)
	assertTagParamError(t, `'key1\`)
}