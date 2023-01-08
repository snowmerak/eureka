package parser

import (
	"fmt"

	"github.com/snowmerak/eureka/token"
)

// Field is a field of a struct.
// <name>: <type>
type Field struct {
	Name string
	Kind int
	Type any
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

	field := &Field{
		Name: string(data[0].Value),
	}
	remains := data

	switch data[2].Kind {
	case token.KindType:
		field.Kind = KindType
		field.Type = string(data[2].Value)
		remains = data[3:]
	case token.KindKeyword:
		switch string(data[2].Value) {
		case "[":
			arr, rem, err := ParseArray(data[2:])
			if err != nil {
				return nil, nil, err
			}
			field.Kind = KindArray
			field.Type = *arr
			remains = rem
		case "struct":
			str, rem, err := ParseStrcut(data[2:])
			if err != nil {
				return nil, nil, err
			}
			field.Kind = KindStruct
			field.Type = *str
			remains = rem
		case "interface":
			ifc, rem, err := ParseInterface(data[2:])
			if err != nil {
				return nil, nil, err
			}
			field.Kind = KindInterface
			field.Type = *ifc
			remains = rem
		default:
			return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
		}
	}

	return field, remains, nil
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
