package main

import (
	"fmt"

	"github.com/snowmerak/eureka/token"
)

func main() {
	a := []byte(`abc + def+789`)

	tokens, err := token.Parse(a)
	if err != nil {
		panic(err)
	}

	for _, t := range tokens {
		fmt.Printf("%s, %d\n", t.Value, t.Kind)
	}
}
