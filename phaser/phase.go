package phaser

import (
	"fmt"

	"github.com/AurelienS/cigare/model"
)

var (
	step = 30
)

func Phase(flight model.Flight) []model.Phase {
	phases := detect(flight)
	phases = merge(phases)
	addDuration(&phases)
	fmt.Println("file: phase.go ~ line 14 ~ phases : ", phases)
	return phases
}

func addDuration(phases *[]model.Phase) {
	for i, p := range *phases {
		(*phases)[i].Duration = p.EndTime.Sub(p.StartTime)
	}
}

func merge(phases []model.Phase) []model.Phase {
	// Check if the input slice is empty
	if len(phases) == 0 {
		return nil
	}

	mergedPhases := make([]model.Phase, 0)
	mergedPhases = append(mergedPhases, phases[0])

	for i := 1; i < len(phases); i++ {
		currentPhase := phases[i]
		lastMergePhase := &mergedPhases[len(mergedPhases)-1]

		// Check if there's a merged phase and it can be extended
		if lastMergePhase.Type == currentPhase.Type {
			lastMergePhase.EndIndex = i
			lastMergePhase.EndTime = currentPhase.EndTime
		} else {
			// If types are different or no merged phase, add the current phase to the result
			mergedPhases = append(mergedPhases, currentPhase)
		}
	}
	return mergedPhases
}
func detect(flight model.Flight) []model.Phase {
	var phases []model.Phase
	for i := 0; i < len(flight.Points); i += step {
		actualStep := step
		if i+step > len(flight.Points)-1 {
			actualStep = len(flight.Points) - i - 1
		}
		turns := NumberOfTurn(flight, i, actualStep)
		if turns.turnCount > 1 {
			phases = append(phases, model.Phase{
				Type:         model.Circling,
				StartIndex:   i,
				EndIndex:     i + actualStep,
				StartTime:    flight.Points[i].Time,
				EndTime:      flight.Points[i+actualStep].Time,
				NumberOfTurn: turns.turnCount,
			})
		} else {
			phases = append(phases, model.Phase{
				Type:       model.Cruising,
				StartIndex: i,
				EndIndex:   i + actualStep,
				StartTime:  flight.Points[i].Time,
				EndTime:    flight.Points[i+actualStep].Time,
			})
		}
	}
	return phases
}

type Turns struct {
	leftTurnCount  int
	rightTurnCount int
	turnCount      int
}

func NumberOfTurn(flight model.Flight, index int, window int) Turns {
	endWindowIndex := index + window
	if len(flight.Points) < endWindowIndex {
		return Turns{}
	}

	left := 0
	right := 0
	totalBearingChange := 0.0
	for i := index; i < endWindowIndex-1; i++ {
		currentPoint := flight.Points[i]

		// Calculate the bearing change from the previous point
		bearingChange := flight.Points[i+1].Bearing - currentPoint.Bearing
		totalBearingChange += bearingChange

		if totalBearingChange > 360 {
			right++
			// totalBearingChange -= 360
		}
		if totalBearingChange < -360 {
			left++
			// totalBearingChange += 360
		}
	}
	return Turns{
		leftTurnCount:  left,
		rightTurnCount: right,
		turnCount:      left + right,
	}
}
