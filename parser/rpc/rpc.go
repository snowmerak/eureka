package rpc

import (
	"bytes"
	"fmt"

	"github.com/snowmerak/eureka/token"
)

const (
	KindStruct = iota
	KindEnum
	KindInterface
	KindFunc
	KindField
)

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

	if data[2].Kind != token.KindType {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
	}

	return &Field{
		Name: string(data[0].Value),
		Type: string(data[2].Value),
	}, data[3:], nil
}

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

type Identifier string

type Func struct {
	Name    string
	Params  []Field
	Returns []Identifier
}

func ParseFunc(data []*token.Token) (*Func, []*token.Token, error) {
	if len(data) < 4 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if data[0].Kind != token.KindKeyword || string(data[0].Value) != "func" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[0].Value))
	}

	if data[1].Kind != token.KindIdentifier {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[1].Value))
	}

	if data[2].Kind != token.KindKeyword || string(data[2].Value) != "(" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
	}

	name := string(data[1].Value)

	params := []Field(nil)

	data = data[3:]
	for {
		if len(data) == 0 {
			return nil, nil, fmt.Errorf("unexpected end of input")
		}

		if data[0].Kind == token.KindKeyword && string(data[0].Value) == ")" {
			break
		}

		if data[0].Kind == token.KindKeyword && string(data[0].Value) == "," {
			data = data[1:]
			continue
		}

		param, remains, err := ParseField(data)
		if err != nil {
			return nil, nil, err
		}

		params = append(params, *param)

		data = remains
	}

	data = data[1:]

	if len(data) < 2 {
		return &Func{
			Name:    name,
			Params:  params,
			Returns: []Identifier(nil),
		}, data, nil
	}

	if data[0].Kind != token.KindKeyword || string(data[0].Value) != "->" {
		return &Func{
			Name:    name,
			Params:  params,
			Returns: []Identifier(nil),
		}, data, nil
	}

	data = data[1:]

	if len(data) > 0 && data[0].Kind == token.KindKeyword && string(data[0].Value) == "(" {
		data = data[1:]
	}

	returns := []Identifier(nil)
	for {
		if len(data) == 0 {
			break
		}

		if data[0].Kind == token.KindKeyword && string(data[0].Value) == "," {
			data = data[1:]
			continue
		}

		if (data[0].Kind == token.KindKeyword && string(data[0].Value) == ")") || (data[0].Kind != token.KindIdentifier && data[0].Kind != token.KindType) {
			break
		}

		if data[0].Kind != token.KindIdentifier && data[0].Kind != token.KindType {
			return nil, nil, fmt.Errorf("unexpected token: %s", string(data[0].Value))
		}

		returns = append(returns, Identifier(string(data[0].Value)))

		data = data[1:]
	}

	return &Func{
		Name:    name,
		Params:  params,
		Returns: returns,
	}, data, nil
}

type Interface struct {
	Name    string
	Methods []Func
}

func ParseInterface(data []*token.Token) (*Interface, []*token.Token, error) {
	if len(data) < 4 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if data[0].Kind != token.KindKeyword || string(data[0].Value) != "interface" {
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

	methods := []Func(nil)
	for {
		if len(data) == 0 {
			return nil, nil, fmt.Errorf("unexpected end of input")
		}

		if data[0].Kind == token.KindKeyword && string(data[0].Value) == "}" {
			break
		}

		method, remains, err := ParseFunc(data)
		if err != nil {
			return nil, nil, err
		}

		data = remains

		methods = append(methods, *method)
	}

	return &Interface{
		Name:    name,
		Methods: methods,
	}, data[1:], nil
}

type Member string

type Enum struct {
	Name   string
	Type   string
	Values []Member
}

func ParseEnum(data []*token.Token) (*Enum, []*token.Token, error) {
	if len(data) < 6 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if data[0].Kind != token.KindKeyword || string(data[0].Value) != "enum" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[0].Value))
	}

	if data[1].Kind != token.KindIdentifier {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[1].Value))
	}

	if data[2].Kind != token.KindKeyword || string(data[2].Value) != ":" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
	}

	if data[3].Kind != token.KindType || !(bytes.HasPrefix(data[3].Value, []byte{'i'}) || bytes.HasPrefix(data[3].Value, []byte{'u'}) || bytes.HasPrefix(data[3].Value, []byte{'f'})) {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[3].Value))
	}

	if data[4].Kind != token.KindKeyword || string(data[4].Value) != "{" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
	}

	name := string(data[1].Value)
	typ := string(data[3].Value)

	data = data[5:]

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
