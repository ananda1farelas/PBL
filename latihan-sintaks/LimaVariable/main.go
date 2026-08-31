package main

import "fmt"

func main() {
	var name string = "Farrelas"
	var age int = 21
	var IPK float64 = 3.12
	var Aktif bool = true
	matkul := []string{"Bahasa Indonesia", "Kewarganeagaraan", "Bahasa Inggris"}

	fmt.Printf("Name: %s (Tipe: string)\n ", name)
	fmt.Printf("Age: %d (Tipe: int)\n ", age)
	fmt.Printf("IPK: %.2f (Tipe: float64)\n ", IPK)
	fmt.Printf("Aktif: %t (Tipe: bool)\n ", Aktif)
	fmt.Printf("Mata Kuliah: %v (Tipe: slice)\n ", matkul)
}
