package main

import "fmt"

func main() {
	employee := Employee{
		Person : Person{
			Name: "wilson", Age: 26,
		},
		EmployeeId: 1,
	}
	employee.PrintInfo()
}

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeId int
}

func (e Employee) PrintInfo() {
	fmt.Println("EmployeeID:", e.EmployeeId)
	fmt.Println("Person.Name:", e.Person.Name)
	fmt.Println("Person.Age:", e.Person.Age)
}