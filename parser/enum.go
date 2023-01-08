package parser

import (
	"bytes"
	"fmt"

	"github.com/snowmerak/eureka/token"
)

// Member is a member of enum.
// <name>
type Member string

// Enum is a enum.
//
//	enum <name> {
//	  <member>
//	  <member>
//	  ...
//	}
//
// or
//
//	enum <name>: <type> {
//	  <member>
//	  <member>
//	  ...
//	}
type Enum struct {
	Name   string
	Type   string
	Values []Member
}

func ParseEnum(data []*token.Token) (*Enum, []*token.Token, error) {
	if len(data) < 4 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if data[0].Kind != token.KindKeyword || string(data[0].Value) != "enum" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[0].Value))
	}

	if data[1].Kind != token.KindIdentifier {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[1].Value))
	}

	name := string(data[1].Value)

	secondValue := string(data[2].Value)
	if data[2].Kind != token.KindKeyword || (secondValue != ":" && secondValue != "{") {
		return nil, nil, fmt.Errorf("unexpected token: %s", secondValue)
	}

	typ := ""

	switch secondValue {
	case ":":
		if len(data) < 5 {
			return nil, nil, fmt.Errorf("unexpected end of input")
		}

		if data[3].Kind != token.KindType || !(bytes.HasPrefix(data[3].Value, []byte{'i'}) || bytes.HasPrefix(data[3].Value, []byte{'u'}) || bytes.HasPrefix(data[3].Value, []byte{'f'})) {
			return nil, nil, fmt.Errorf("unexpected token: %s", string(data[3].Value))
		}

		if data[4].Kind != token.KindKeyword || string(data[4].Value) != "{" {
			return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
		}

		typ = string(data[3].Value)
		data = data[5:]
	case "{":
		typ = "i32"
		data = data[3:]
	}

	values := []Member(nil)
	for {
		if len(data) == 0 {
			return nil, nil, fmt.Errorf("unexpected end of input")
		}

		if data[0].Kind == token.KindKeyword && string(data[0].Value) == "}" {
			break
		}

		if data[0].Kind != token.KindIdentifier {
			return nil, nil, fmt.Errorf("unexpected token: %s", string(data[0].Value))
		}

		values = append(values, Member(string(data[0].Value)))

		data = data[1:]
	}

	return &Enum{
		Name:   name,
		Type:   typ,
		Values: values,
	}, data[1:], nil
}
