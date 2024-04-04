// values.go
package main

import "fmt"

func main() {

	fmt.Println("go" + "lang")
	fmt.Printf("golang starting...")
	fmt.Println("1+1 =", 1+1)
	fmt.Println("7.0/3.0 =", 7.0/3.0)
	fmt.Println("x=", 9%2)
	fmt.Println("check", 90 == 90.000)
	// fmt.Println("check", 90 === 90.000) // errors
	fmt.Println(true && false)
	fmt.Println(true || false)
	fmt.Println(!true)
}
