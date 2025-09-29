package main

import "fmt"

func main() {
	var s Shape
	r := Rectangle{}
	s = r
	s.Area()
	s.Perimeter()

	c := Circle{}
	s = c
	s.Area()
	s.Perimeter()

}

func (r Rectangle) Area() {
	fmt.Println("Area Rectangle")

}

func (c Circle) Area() {
	fmt.Println("Area Circle")

}

func (r Rectangle) Perimeter() {
	fmt.Println("Perimeter Rectangle")

}

func (c Circle) Perimeter() {
	fmt.Println("Perimeter Circle")

}

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
}

type Circle struct {
}