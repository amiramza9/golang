// lucky_no
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// Get a random number between 0 and 99 inclusive.
	n := rand.Intn(100)

	// Print it out.
	fmt.Printf("Your lucky number is %d!\n", n)
}
