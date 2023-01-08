package parser

import (
	"fmt"

	"github.com/snowmerak/eureka/token"
)

// Field is a field of a struct.
// <name>: <type>
type Field struct {
	Name string
	Type string
}

func ParseField(data []*token.Token) (*Field, []*token.Token, error) {
	if len(data) < 3 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if data[0].Kind != token.KindIdentifier {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[0].Value))
	}

	if data[1].Kind != token.KindKeyword || string(data[1].Value) != ":" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[1].Value))
	}

	if data[2].Kind != token.KindType && data[2].Kind != token.KindIdentifier {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
	}

	return &Field{
		Name: string(data[0].Value),
		Type: string(data[2].Value),
	}, data[3:], nil
}

// struct is a structure.
//
//	struct <name> {
//		<field>
//		<field>
//		...
//	}
type Struct struct {
	Name   string
	Fields []Field
}

func ParseStrcut(data []*token.Token) (*Struct, []*token.Token, error) {
	if len(data) < 4 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if data[0].Kind != token.KindKeyword || string(data[0].Value) != "struct" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[0].Value))
	}

	if data[1].Kind != token.KindIdentifier {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[1].Value))
	}

	if data[2].Kind != token.KindKeyword || string(data[2].Value) != "{" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
	}

	name := string(data[1].Value)

	data = data[3:]

	fields := []Field(nil)
	for {
		if len(data) == 0 {
			return nil, nil, fmt.Errorf("unexpected end of input")
		}

		if data[0].Kind == token.KindKeyword && string(data[0].Value) == "}" {
			break
		}

		field, remains, err := ParseField(data)
		if err != nil {
			return nil, nil, err
		}

		data = remains

		fields = append(fields, *field)
	}

	return &Struct{
		Name:   name,
		Fields: fields,
	}, data[1:], nil
}
