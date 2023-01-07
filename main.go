package main

import (
	"fmt"
	"log"

	"github.com/snowmerak/eureka/token"
)

func main() {
	a := []byte(`select * from table where id = 1 and name = "snowmerak"`)
	fmt.Println([]byte(` = 1 and name = "snowmerak"`))

	tokens, err := token.Parse(a)
	if err != nil {
		log.Println(err)
	}

	for _, t := range tokens {
		fmt.Printf("%s, %d\n", t.Value, t.Kind)
	}
}
