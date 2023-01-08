package parser

import (
	"fmt"
	"strconv"

	"github.com/snowmerak/eureka/token"
)

type Array struct {
	Length int
	Type   string
	Values []string
}

func ParseArray(data []*token.Token) (*Array, []*token.Token, error) {
	if len(data) < 4 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if data[0].Kind != token.KindKeyword || string(data[0].Value) != "[" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[0].Value))
	}

	if data[1].Kind != token.KindInteger {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[1].Value))
	}

	if data[2].Kind != token.KindKeyword || string(data[2].Value) != "]" {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[2].Value))
	}

	length, err := strconv.Atoi(string(data[1].Value))
	if err != nil {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[1].Value))
	}

	if data[3].Kind != token.KindType {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(data[3].Value))
	}

	typ := string(data[3].Value)

	if len(data) < 5 || data[4].Kind != token.KindKeyword || string(data[4].Value) != "{" {
		return &Array{
			Length: length,
			Type:   typ,
			Values: []string{},
		}, data[3:], nil
	}

	var values []string
	for i := 5; i < len(data); i++ {
		if data[i].Kind == token.KindKeyword && string(data[i].Value) == "," {
			continue
		}

		if data[i].Kind == token.KindKeyword && string(data[i].Value) == "}" {
			if len(values) > length {
				return nil, nil, fmt.Errorf("unexpected token: %s", string(data[i].Value))
			}
			return &Array{
				Length: length,
				Type:   typ,
				Values: values,
			}, data[i+1:], nil
		}

		if data[i].Kind != token.KindString && data[i].Kind != token.KindInteger && data[i].Kind != token.KindFloat && data[i].Kind != token.KindBool && data[i].Kind != token.KindIdentifier {
			return nil, nil, fmt.Errorf("unexpected token: %s", string(data[i].Value))
		}

		values = append(values, string(data[i].Value))
	}

	return nil, nil, fmt.Errorf("unexpected end of input")
}
