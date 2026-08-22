package main

import "fmt"

func main() {
	var numbA int = 10
	var numbB *int = &numbA
	numbC := numbB
	var numbD **int = &numbB

	fmt.Println("Nilai dari numbA:", numbA)
	fmt.Println("Nilai dari numbB:", *numbB)
	fmt.Println("Nilai dari numbC:", *numbC)
	fmt.Println("Nilai dari numbD:", **numbD)
}
