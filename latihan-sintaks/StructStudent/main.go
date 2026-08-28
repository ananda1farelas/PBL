package main

import (
	"fmt"
)

type Student struct {
	ID       int
	Name     string
	Grade    float64
	IsActive bool
}

func (s Student) GetInfo() string {
	status := "Tidak Aktif"
	if s.IsActive {
		status = "Aktif"
	}
	return fmt.Sprintf("ID: %d, Nama: %s, Nilai: %.2f, Status: %s", s.ID, s.Name, s.Grade, status)
}

func (s *Student) UpdateGrade(grade float64) {
	s.Grade = grade
}

func (s *Student) Activate() {
	s.IsActive = true
}

func (s *Student) Deactivate() {
	s.IsActive = false
}

func main() {
	student1 := Student{
		ID:       937538204,
		Name:     "Farrelas",
		Grade:    3.12,
		IsActive: true,
	}

	fmt.Println("Informasi Mahasiswa:")
	fmt.Println(student1.GetInfo())

	student1.Activate()
	student1.UpdateGrade(3.34)

	fmt.Println("\nSetelah Update:")
	fmt.Println(student1.GetInfo())

	student1.Deactivate()

	fmt.Println("\nSetelah Dinonaktifkan:")
	fmt.Println(student1.GetInfo())
}
