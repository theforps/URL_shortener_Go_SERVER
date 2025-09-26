package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {

	//test := "https://google.com"
	//base := "http://127.0.0.1/"

	config, _ := Configuration()

	fmt.Println(DbCreate())

	fmt.Printf("%s%s", config.Domain, GenerateString(7))

}

func GenerateString(length int) string {

	result := ""

	rangeSymbols := []rune("QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890")
	for i := 0; i < length; i++ {
		randIndex := rand.IntN(len(rangeSymbols))
		result += string(rangeSymbols[randIndex])
	}

	return result
}
