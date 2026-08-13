package main

import "fmt"

func main() {
	dataMahasiswa := make(map[string]float64)

	dataMahasiswa["Farrelas"] = 3.12
	dataMahasiswa["Ditya"] = 3.45
	dataMahasiswa["Hanum"] = 3.67
	fmt.Println("Berhasil Menambahkan Data Mahasiswa")

	//Cek data
	if nilai, exists := dataMahasiswa["Farrelas"]; exists {
		fmt.Printf("Data Farrelas ditemukan! Nilai: %.2f\n", nilai)
	} else {
		fmt.Println("Data Farrelas tidak ditemukan.")
	}

	if nilai, exists := dataMahasiswa["Kalisa"]; exists {
		fmt.Printf("Data Kalisa ditemukan! Nilai: %.2f\n", nilai)
	} else {
		fmt.Println("Data Kalisa tidak ditemukan dalam map.")
	}

	//Menelusuri Map
	for namaMhs, nilaiMhs := range dataMahasiswa {
		fmt.Printf("Mahasiswa: %-8s | Nilai: %.2f\n", namaMhs, nilaiMhs)
	}

	//Menghapus data dari map
	delete(dataMahasiswa, "Ditya")
	fmt.Println("[-] Data 'Ditya' telah dihapus dari map.")

	//Menelusuri kembali buat ngecek apa kah data uda dihapus
	for namaMhs, nilaiMhs := range dataMahasiswa {
		fmt.Printf("Mahasiswa: %-8s | Nilai: %.2f\n", namaMhs, nilaiMhs)
	}
}
