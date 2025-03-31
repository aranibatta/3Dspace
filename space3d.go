package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
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

// LoadPointsFromCSV loads 3D points from a CSV file
// The CSV should have at least 3 columns for X, Y, Z coordinates
func (s *Space3D) LoadPointsFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header if exists (optional)
	// If your CSV has no header, comment this line out
	_, err = reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading CSV record: %w", err)
		}

		// Ensure we have at least 3 columns for X, Y, Z
		if len(record) < 3 {
			return fmt.Errorf("invalid CSV format: need at least 3 columns for X, Y, Z coordinates")
		}

		// Parse the coordinates
		x, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return fmt.Errorf("invalid X coordinate: %w", err)
		}

		y, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return fmt.Errorf("invalid Y coordinate: %w", err)
		}

		z, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return fmt.Errorf("invalid Z coordinate: %w", err)
		}

		// Add the point to our space
		s.AddPoint(NewPoint3D(x, y, z))
	}

	return nil
}

// SavePointsToCSV saves all points in the space to a CSV file
func (s *Space3D) SavePointsToCSV(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"X", "Y", "Z"})
	if err != nil {
		return fmt.Errorf("error writing CSV header: %w", err)
	}

	// Write points
	for _, point := range s.Points {
		record := []string{
			strconv.FormatFloat(point.X, 'f', -1, 64),
			strconv.FormatFloat(point.Y, 'f', -1, 64),
			strconv.FormatFloat(point.Z, 'f', -1, 64),
		}
		err := writer.Write(record)
		if err != nil {
			return fmt.Errorf("error writing point to CSV: %w", err)
		}
	}

	return nil
}
