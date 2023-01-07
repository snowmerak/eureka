package token

import "fmt"

const (
	KindKeyworkd = iota
	KindString
	KindCharactor
	KindInteger
	KindFloat
	KindBool
	KindSymbol
	KindIdentifier
	KindSpace
)

var Keywords = [][]byte{
	[]byte("func"),
	[]byte("return"),
	[]byte("if"),
	[]byte("for"),
	[]byte("break"),
	[]byte("continue"),
	[]byte("let"),
	[]byte("mut"),
	[]byte("type"),
	[]byte("struct"),
	[]byte("enum"),
	[]byte("interface"),
	[]byte("match"),
	[]byte("=>"),
	[]byte("("),
	[]byte(")"),
	[]byte("{"),
	[]byte("}"),
	[]byte(":"),
	[]byte(";"),
	[]byte(","),
	[]byte("->"),
	[]byte("<-"),
	[]byte("<->"),
	[]byte("."),
	[]byte("="),
	[]byte("=="),
	[]byte("!="),
	[]byte("<"),
	[]byte(">"),
	[]byte("<="),
	[]byte(">="),
	[]byte("+"),
	[]byte("-"),
	[]byte("*"),
	[]byte("/"),
	[]byte("%"),
	[]byte("&&"),
	[]byte("||"),
	[]byte("!"),
	[]byte("&"),
	[]byte("|"),
	[]byte("^"),
	[]byte("<<"),
	[]byte(">>"),
	[]byte("  "),
}

func ParseKeyword(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	for _, v := range Keywords {
		if len(buf) < len(v) {
			continue
		}

		if string(buf[:len(v)]) == string(v) {
			return &Token{
				Value: v,
				Kind:  KindKeyworkd,
			}, buf[len(v):], nil
		}
	}

	return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
}

func ParseString(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if buf[0] != '"' {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	buf = buf[1:]
	value := []byte(nil)
	for i := 0; i < len(buf); i++ {
		if buf[i] == '"' {
			return &Token{
				Value: value,
				Kind:  KindString,
			}, buf[i+1:], nil
		}

		value = append(value, buf[i])
	}

	return nil, nil, fmt.Errorf("unexpected end of input")
}

func ParseCharactor(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if buf[0] != '\'' {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	buf = buf[1:]
	value := []byte(nil)
	for i := 0; i < len(buf); i++ {
		if buf[i] == '\'' {
			return &Token{
				Value: value,
				Kind:  KindCharactor,
			}, buf[i+1:], nil
		}

		value = append(value, buf[i])
	}

	return nil, nil, fmt.Errorf("unexpected end of input")
}

func ParseInteger(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if buf[0] < '0' || buf[0] > '9' {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	value := []byte(nil)
	for i := 0; i < len(buf); i++ {
		if buf[i] < '0' || buf[i] > '9' {
			return &Token{
				Value: value,
				Kind:  KindInteger,
			}, buf[i:], nil
		}

		value = append(value, buf[i])
	}

	return &Token{
		Value: value,
		Kind:  KindInteger,
	}, nil, nil
}

func ParseFloat(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if buf[0] < '0' || buf[0] > '9' {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	dotExists := false
	value := []byte(nil)
	for i := 0; i < len(buf); i++ {
		if buf[i] < '0' || buf[i] > '9' {
			if buf[i] == '.' {
				if dotExists {
					return &Token{
						Value: value,
						Kind:  KindFloat,
					}, buf[i:], nil
				}
				value = append(value, buf[i])
				dotExists = true
				continue
			}

			return &Token{
				Value: value,
				Kind:  KindFloat,
			}, buf[i:], nil
		}

		value = append(value, buf[i])
	}

	return &Token{
		Value: value,
		Kind:  KindFloat,
	}, nil, nil
}

func ParseBool(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if buf[0] != 't' && buf[0] != 'f' {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	if len(buf) < 4 {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	if string(buf[:4]) == "true" {
		return &Token{
			Value: []byte("true"),
			Kind:  KindBool,
		}, buf[4:], nil
	}

	if len(buf) < 5 {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	if string(buf[:5]) == "false" {
		return &Token{
			Value: []byte("false"),
			Kind:  KindBool,
		}, buf[5:], nil
	}

	return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
}

func ParseIdentifier(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if buf[0] < 'a' || buf[0] > 'z' {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	value := []byte(nil)
	for i := 0; i < len(buf); i++ {
		if (buf[i] < 'a' || buf[i] > 'z') && (buf[i] < '0' || buf[i] > '9') && buf[i] != '_' {
			return &Token{
				Value: value,
				Kind:  KindIdentifier,
			}, buf[i:], nil
		}

		value = append(value, buf[i])
	}

	return &Token{
		Value: value,
		Kind:  KindIdentifier,
	}, nil, nil
}

func ParseSymbol(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	if buf[0] != ':' {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	if len(buf) > 2 && (buf[1] < 'a' || buf[1] > 'z') {
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	buf = buf[1:]

	value := []byte(":")
	for i := 0; i < len(buf); i++ {
		if (buf[i] < 'a' || buf[i] > 'z') && (buf[i] < '0' || buf[i] > '9') {
			return &Token{
				Value: value,
				Kind:  KindSymbol,
			}, buf[len(value)-1:], nil
		}

		value = append(value, buf[i])
	}

	return &Token{
		Value: value,
		Kind:  KindSymbol,
	}, buf[len(value)-1:], nil
}

func ParseSpace(buf []byte) (*Token, []byte, error) {
	if len(buf) == 0 {
		return nil, nil, fmt.Errorf("unexpected end of input")
	}

	switch buf[0] {
	case ' ', '\t', '\r', '\n':
	default:
		return nil, nil, fmt.Errorf("unexpected token: %s", string(buf))
	}

	return &Token{
		Value: buf[:1],
		Kind:  KindSpace,
	}, buf[1:], nil
}

var parsers = []func([]byte) (*Token, []byte, error){
	ParseKeyword,
	ParseString,
	ParseCharactor,
	ParseFloat,
	ParseInteger,
	ParseBool,
	ParseSymbol,
	ParseIdentifier,
	ParseSpace,
}

func Parse(buf []byte) ([]*Token, error) {
	tokens := []*Token(nil)

	globalError := error(nil)
	for len(buf) > 0 {
		for _, parser := range parsers {
			token, rest, err := parser(buf)
			if err != nil {
				globalError = err
				continue
			}
			buf = rest
			tokens = append(tokens, token)
			globalError = nil
			break
		}

		if globalError != nil {
			return tokens, globalError
		}
	}

	return tokens, nil
}
