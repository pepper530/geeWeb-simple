package main

import "fmt"

func main() {
	word := "pepper"
	for _, ch := range word {
		fmt.Println("ch=", ch)
		ch -= 'a'
		fmt.Println("ch-a=", ch)
	}
}
