package main

import (
	"fmt"
	"math"
)

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

// function to calculate the distance between two points
// in a 2D plane
func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func printTheString(x string) {
	fmt.Println(x)
}

// A function that calculates the distance between two points in a 3D plane
func distance3D(x1, y1, z1, x2, y2, z2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2) + math.Pow(z2-z1, 2))
}

func main() {
	fmt.Println("Hello World")
	
	// Create a 3D space
	space := NewSpace3D()
	
	// Add some points
	p1 := NewPoint3D(0, 0, 0)
	p2 := NewPoint3D(3, 4, 0)
	p3 := NewPoint3D(3, 4, 5)
	
	space.AddPoint(p1)
	space.AddPoint(p2)
	space.AddPoint(p3)
	
	// Calculate and print distances
	fmt.Printf("Distance p1 to p2: %.2f\n", Distance(p1, p2))
	fmt.Printf("Distance p2 to p3: %.2f\n", p2.DistanceTo(p3))
	fmt.Printf("Manhattan distance p1 to p3: %.2f\n", ManhattanDistance(p1, p3))
}
