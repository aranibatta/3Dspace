package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
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

func generateSampleCSV(fileName string) error {
	space := NewSpace3D()

	// Add some sample points
	space.AddPoint(NewPoint3D(0, 0, 0))
	space.AddPoint(NewPoint3D(1, 0, 0))
	space.AddPoint(NewPoint3D(0, 1, 0))
	space.AddPoint(NewPoint3D(0, 0, 1))
	space.AddPoint(NewPoint3D(1, 1, 0))
	space.AddPoint(NewPoint3D(1, 0, 1))
	space.AddPoint(NewPoint3D(0, 1, 1))
	space.AddPoint(NewPoint3D(1, 1, 1))

	// Generate some points on a helix
	for t := 0.0; t < 10; t += 0.1 {
		x := math.Cos(t)
		y := math.Sin(t)
		z := t / 3
		space.AddPoint(NewPoint3D(x, y, z))
	}

	return space.SavePointsToCSV(fileName)
}

func main() {
	// Command line flags
	csvFile := flag.String("csv", "", "Path to CSV file with 3D points")
	generateSample := flag.String("generate", "", "Generate a sample CSV file at the specified path")
	functionStr := flag.String("function", "", "Mathematical function to visualize (e.g., 'sin(x) * cos(y)')")
	xMin := flag.Float64("xmin", -5.0, "Minimum x value for function visualization")
	xMax := flag.Float64("xmax", 5.0, "Maximum x value for function visualization")
	yMin := flag.Float64("ymin", -5.0, "Minimum y value for function visualization")
	yMax := flag.Float64("ymax", 5.0, "Maximum y value for function visualization")
	step := flag.Float64("step", 0.2, "Step size for function visualization")

	flag.Parse()

	// Generate sample data if requested
	if *generateSample != "" {
		fmt.Printf("Generating sample CSV file at %s\n", *generateSample)
		if err := generateSampleCSV(*generateSample); err != nil {
			log.Fatalf("Failed to generate sample CSV: %v", err)
		}
		fmt.Println("Sample CSV file generated successfully")
		os.Exit(0)
	}

	// Create a 3D space
	space := NewSpace3D()

	// Check if function visualization is requested
	if *functionStr != "" {
		fmt.Printf("Generating points from function: %s\n", *functionStr)
		fmt.Printf("Range: x=[%.2f, %.2f], y=[%.2f, %.2f], step=%.2f\n", *xMin, *xMax, *yMin, *yMax, *step)
		
		if err := GeneratePointsFromFunction(space, *functionStr, *xMin, *xMax, *yMin, *yMax, *step); err != nil {
			log.Fatalf("Error generating points from function: %v", err)
		}
		
		fmt.Printf("Generated %d points from function\n", len(space.Points))
	// Load points from CSV if provided
	} else if *csvFile != "" {
		fmt.Printf("Loading points from CSV file: %s\n", *csvFile)
		if err := space.LoadPointsFromCSV(*csvFile); err != nil {
			log.Fatalf("Error loading CSV file: %v", err)
		}
		fmt.Printf("Loaded %d points from CSV\n", len(space.Points))
	} else {
		// Add some default points if no CSV provided
		fmt.Println("No function or CSV file specified, using default points")
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

	// Create and run the visualizer
	visualizer := NewVisualizer(space)
	visualizer.Run()
}
