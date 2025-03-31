# 3D Points Visualizer

A simple Go application for visualizing 3D points from CSV files.

## Features

- Import 3D points from CSV files
- Generate sample 3D point data (including a helix)
- Interactive 3D visualization with rotation and scaling
- Calculate distances between points (Euclidean and Manhattan)

## Requirements

- Go 1.20 or higher
- Fyne GUI toolkit (automatically installed via go modules)

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd learnGo

# Install dependencies
go mod tidy
```

## Usage

### Running with Default Points

```bash
go run .
```

### Generating a Sample CSV File

```bash
go run . -generate sample_points.csv
```

### Loading Points from a CSV File

```bash
go run . -csv your_points.csv
```

## CSV File Format

The CSV file should have at least 3 columns for X, Y, and Z coordinates. The first row should be a header row.

Example:
```
X,Y,Z
0,0,0
1,0,0
0,1,0
0,0,1
```

## Controls

In the visualization window:
- Use sliders to rotate the view around X, Y, and Z axes
- Use the scale slider to zoom in and out

## Building

```bash
go build
```

This will create an executable file that you can run directly.