package main

import "fmt"

func swap(a, b *int) {
	*a, *b = *b, *a //menggunakan pointer buat nuker variable a dan b
}

func updateSlice(s *[]string, newValue string) {
	*s = append(*s, newValue) //menggunakan pointer buat  nambahin data baru ke slice
}

func swapValue(a, b int) {
	a, b = b, a //ini hanya nuker niali di local function, hasil nya akan tetap sama seperti yang di input
}

func updateSliceValue(s []string, newValue string) {
	s = append(s, newValue) //ini nambah data di local function juga, jadi hasilnya bakal tetep sama seperti yang di input
}

func main() {
	//demo swap dan updateslice pake pointer
	x, y := 15, 30
	fmt.Printf("sebelum swap: x = %d, y = %d\n", x, y)
	swap(&x, &y) //operasi swap
	fmt.Printf("sesudah swap: x = %d, y = %d\n", x, y)

	hobi := []string{"Membaca", "Menulis"}
	fmt.Printf("sebelum UpdateSlice: %v (Jumlah data: %d)\n", hobi, len(hobi))
	updateSlice(&hobi, "Membangun restoran di roblok")
	fmt.Printf("Sesudah di UpdateSlice : %v (Jumlah data: %d),\n", hobi, len(hobi))

	//pembuktian menggunakan value
	x, y = 12, 24
	fmt.Printf("sebelum swapValue: x = %d, y = %d\n", x, y)
	swapValue(x, y)
	fmt.Printf("sesudah swapValue: x = %d, y = %d\n", x, y)

	hobi = []string{"Menulis", "Memasak"}
	fmt.Printf("sebelum UpdateSliceValue: %v (Jumlah data: %d)\n", hobi, len(hobi))
	updateSliceValue(hobi, "Membangun restoran di roblok")
	fmt.Printf("Sesudah di UpdateSliceValue : %v (Jumlah data: %d),\n", hobi, len(hobi))
}
