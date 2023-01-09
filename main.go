package main

import (
	"fmt"
	"log"
	"os"

	"github.com/snowmerak/eureka/parser"
	"github.com/snowmerak/eureka/token"
)

func main() {
	data, err := os.ReadFile("test.csc")
	if err != nil {
		panic(err)
	}

	tokens, err := token.ParseWithoutSpace(data)
	if err != nil {
		log.Println(err)
	}

	person, remains, err := parser.ParseStrcut(tokens)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(person)

	itf, remains, err := parser.ParseInterface(remains)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(itf)

	enum, remains, err := parser.ParseEnum(remains)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(enum)

	arr, remains, err := parser.ParseArray(remains)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(arr)
}
