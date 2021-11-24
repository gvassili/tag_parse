package tag_parser

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type parserState int

const (
	isQuoted  parserState = 1 << iota
	isEscaped
	isValue
)

type Param struct {
	Key    string
	Values []string
}

// Parse take a struct tag string and return a slice of Param.
// The key need to be comma separated, key value(s) need to be preceded by an equal character and can be space separated to add multiple values.
// Space and commas can be escaped with in single quoted block, single quote and backslash can be escaped with a backslash.
// Leading and trailing space are ignored.
// An empty tag is valid and will return a nil slice of Param.
// If tag contain multiple keys on a param, an unterminated escape or quote or an invalid escape sequence, Tag will return a nil slice of Param and a non-null error.
func Parse(tag string) ([]Param, error) {
	r := strings.NewReader(tag)
	sb := strings.Builder{}
	state := parserState(0)
	param := Param{}
	params := make([]Param, 0)
	for {
		b, err := r.ReadByte()
		pushByte := func() {
			sb.WriteByte(b)
		}
		pushWord := func() error {
			if state & isValue == 0 {
				if sb.String() != "" {
					if param.Key != "" {
						return errors.New("can only have one key")
					} else if sb.String() == "" {
						return errors.New("key is empty")
					}
					param.Key = sb.String()
				}
			} else if sb.String() != "" {
				param.Values = append(param.Values, sb.String())
			}
			sb.Reset()
			return nil
		}

		if err != nil {
			if err == io.EOF {
				if sb.String() != "" {
					if err := pushWord(); err != nil {
						return nil, err
					}
				}
				break
			} else {
				return nil, err
			}
		}

		if state & isEscaped != 0 {
			state &= ^isEscaped
			switch b {
			case '\\':
				pushByte()
			case '\'':
				pushByte()
			default:
				return nil, fmt.Errorf("'\\%c' isn't a valid escape sequence", b)
			}
			continue
		}

		switch b {
		case '\\':
			state |= isEscaped
		case '\'':
			state ^= isQuoted
		case ' ':
			if state & isQuoted == 0 {
				if err := pushWord(); err != nil {
					return nil, err
				}
			} else {
				pushByte()
			}
		case ',':
			if state & isQuoted == 0 {
				if sb.String() != "" {
					if err := pushWord(); err != nil {
						return nil, err
					}
				}
				if param.Key == "" {
					return nil, errors.New("key is empty")
				}
				state = 0
				params = append(params, param)
				param = Param{}
			} else {
				pushByte()
			}
		case '=':
			if err := pushWord(); err != nil {
				return nil, err
			}
			state |= isValue
		default:
			pushByte()
		}
	}
	if param.Key != "" {
		params = append(params, param)
		param = Param{}
	}
	if state & isQuoted != 0 {
		return nil, errors.New("unterminated quote")
	}
	if state & isEscaped != 0 {
		return nil, errors.New("unterminated escape sequence")
	}
	return params, nil
}