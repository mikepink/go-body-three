package main

import (
	"fmt"
	"math"
	"sync"
)

type Frame struct {
	Ids       []uint16  `json:"ids"`
	Positions []float64 `json:"positions"`
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Point struct {
	Acceleration *Vector3
	Id           uint16
	Mass         float64
	Position     *Vector3
	Velocity     *Vector3
}

type NBodySimulation struct {
	Points []*Point
}

func (n *NBodySimulation) AddPoint(p *Point) []*Point {
	p.Id = uint16(len(n.Points) + 1)
	n.Points = append(n.Points, p)
	return n.Points
}

const STEP_SIZE = 1.0 / 60.0
const G = 0.0000000000667

func Simulator(wg sync.WaitGroup, positionChan chan<- *Frame, endSimulationChan chan<- bool) {
	defer wg.Done()
	fmt.Println("Simulator called")
	n := NBodySimulation{
		Points: make([]*Point, 0),
	}
	n.AddPoint(&Point{
		Acceleration: &Vector3{0, 0, 0},
		Mass:         500000000,
		Position:     &Vector3{0, 0, 0},
		Velocity:     &Vector3{0, 0, 0},
	})
	n.AddPoint(&Point{
		Acceleration: &Vector3{0, 0, 0},
		Mass:         2,
		Position:     &Vector3{5, 0, 0},
		Velocity:     &Vector3{0, -0.03, 0},
	})
	n.AddPoint(&Point{
		Acceleration: &Vector3{0, 0, 0},
		Mass:         10000,
		Position:     &Vector3{-12, 0, 0},
		Velocity:     &Vector3{0, 0.015, 0},
	})
	n.AddPoint(&Point{
		Acceleration: &Vector3{0, 0, 0},
		Mass:         1000000,
		Position:     &Vector3{10, 0, -3},
		Velocity:     &Vector3{0, -0.01, -0.01},
	})
	n.AddPoint(&Point{
		Acceleration: &Vector3{0, 0, 0},
		Mass:         10000000,
		Position:     &Vector3{20, 20, -12},
		Velocity:     &Vector3{0, 0, 0.01},
	})
	fmt.Println("Created new simulator struct")
	dT := 1.0
	for i := 0; i < 1500000; i++ {
		numPoints := len(n.Points)
		idBuffer := make([]uint16, numPoints)
		positionsBuffer := make([]float64, numPoints*3)
		for i1 := 0; i1 < numPoints; i1++ {
			for i2 := 0; i2 < numPoints; i2++ {
				// Don't need to calculate effect of force on itself
				if i1 == i2 {
					continue
				}

				n.Points[i1].Velocity.X += n.Points[i1].Acceleration.X * dT * 0.5
				n.Points[i1].Velocity.Y += n.Points[i1].Acceleration.Y * dT * 0.5
				n.Points[i1].Velocity.Z += n.Points[i1].Acceleration.Z * dT * 0.5

				n.Points[i1].Position.X += n.Points[i1].Velocity.X * dT
				n.Points[i1].Position.Y += n.Points[i1].Velocity.Y * dT
				n.Points[i1].Position.Z += n.Points[i1].Velocity.Z * dT

				dx := n.Points[i2].Position.X - n.Points[i1].Position.X
				dy := n.Points[i2].Position.Y - n.Points[i1].Position.Y
				dz := n.Points[i2].Position.Z - n.Points[i1].Position.Z

				invR3 := math.Pow(
					math.Pow(dx, 2)+math.Pow(dy, 2)+math.Pow(dz, 2),
					-1.5,
				)
				coeff := invR3 * n.Points[i2].Mass * G

				n.Points[i1].Acceleration.X = dx * coeff
				n.Points[i1].Acceleration.Y = dy * coeff
				n.Points[i1].Acceleration.Z = dz * coeff

				n.Points[i1].Velocity.X += n.Points[i1].Acceleration.X * dT * 0.5
				n.Points[i1].Velocity.Y += n.Points[i1].Acceleration.Y * dT * 0.5
				n.Points[i1].Velocity.Z += n.Points[i1].Acceleration.Z * dT * 0.5
			}

			positionsBuffer[i1*3] = n.Points[i1].Position.X
			positionsBuffer[i1*3+1] = n.Points[i1].Position.Y
			positionsBuffer[i1*3+2] = n.Points[i1].Position.Z
			idBuffer[i1] = n.Points[i1].Id
		}

		positionChan <- &Frame{
			Ids:       idBuffer,
			Positions: positionsBuffer,
		}
	}
	endSimulationChan <- true
}
