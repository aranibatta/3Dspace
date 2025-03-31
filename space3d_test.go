package main

import (
	"math"
	"testing"
)

func TestPoint3D(t *testing.T) {
	p := NewPoint3D(1, 2, 3)

	if p.X != 1 || p.Y != 2 || p.Z != 3 {
		t.Errorf("Point3D not created correctly, got %v", p)
	}
}

func TestDistance(t *testing.T) {
	p1 := NewPoint3D(0, 0, 0)
	p2 := NewPoint3D(1, 1, 1)

	expected := math.Sqrt(3)
	result := Distance(p1, p2)

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Distance calculation incorrect, expected %v, got %v", expected, result)
	}

	// Test method version
	result = p1.DistanceTo(p2)
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Point method DistanceTo incorrect, expected %v, got %v", expected, result)
	}
}

func TestManhattanDistance(t *testing.T) {
	p1 := NewPoint3D(0, 0, 0)
	p2 := NewPoint3D(1, 1, 1)

	expected := 3.0
	result := ManhattanDistance(p1, p2)

	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("Manhattan distance calculation incorrect, expected %v, got %v", expected, result)
	}
}

func TestSpace3D(t *testing.T) {
	space := NewSpace3D()

	p1 := NewPoint3D(1, 2, 3)
	p2 := NewPoint3D(4, 5, 6)

	space.AddPoint(p1)
	space.AddPoint(p2)

	if len(space.Points) != 2 {
		t.Errorf("Expected 2 points in space, got %d", len(space.Points))
	}

	if space.Points[0] != p1 || space.Points[1] != p2 {
		t.Errorf("Points not stored correctly in space")
	}
}
