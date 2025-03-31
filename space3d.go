package main

import (
	"math"
)

// Point3D represents a point in 3D space
type Point3D struct {
	X, Y, Z float64
}

// NewPoint3D creates a new 3D point
func NewPoint3D(x, y, z float64) Point3D {
	return Point3D{X: x, Y: y, Z: z}
}

// Space3D represents a collection of points in 3D space
type Space3D struct {
	Points []Point3D
}

// NewSpace3D creates a new empty 3D space
func NewSpace3D() *Space3D {
	return &Space3D{
		Points: make([]Point3D, 0),
	}
}

// AddPoint adds a point to the 3D space
func (s *Space3D) AddPoint(p Point3D) {
	s.Points = append(s.Points, p)
}

// Distance calculates the Euclidean distance between two 3D points
func Distance(p1, p2 Point3D) float64 {
	return math.Sqrt(
		math.Pow(p2.X-p1.X, 2) +
			math.Pow(p2.Y-p1.Y, 2) +
			math.Pow(p2.Z-p1.Z, 2))
}

// ManhattanDistance calculates the Manhattan distance between two 3D points
func ManhattanDistance(p1, p2 Point3D) float64 {
	return math.Abs(p2.X-p1.X) + math.Abs(p2.Y-p1.Y) + math.Abs(p2.Z-p1.Z)
}

// Method version of Distance calculation
func (p1 Point3D) DistanceTo(p2 Point3D) float64 {
	return Distance(p1, p2)
}