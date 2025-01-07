package main
import (
	"fmt"
)

func maintyt() {
	var fName string
	fmt.Print("Enter your first name: ")
	fmt.Scanln(&fName)

	var lName string
	fmt.Print("Enter your last name: ")
	fmt.Scanln(&lName)

	fmt.Println("Hello, " + fName + " " + lName)
}