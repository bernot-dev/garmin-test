package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Daily is a summary of daily activity
type Daily struct {
	UserID                             string
	UserAccessToken                    string
	SummaryID                          string
	CalendarDate                       string
	ActivityType                       string
	ActiveKilocalories                 uint
	BMRKilocalories                    uint
	ConsumedCalories                   uint
	Steps                              uint
	DistanceInMeters                   float32
	DurationInSeconds                  uint
	ActiveTimeInSeconds                uint
	StartTimeInSeconds                 uint
	StartTimeOffsetInSeconds           uint
	ModerateIntensityDurationInSeconds uint
	VigorousIntensityDurationInSeconds uint
	FloorsClimbed                      uint
	MinHeartRateInBeatsPerMinute       uint
	AverageHeartRateInBeatsPerMinute   uint
	MaxHeartRateInBeatsPerMinute       uint
	AverageStressLevel                 uint
	MaxStressLevel                     uint
	StressDurationInSeconds            uint
	RestStressDurationInSeconds        uint
	ActivityStressDurationInSeconds    uint
	LowStressDurationInSeconds         uint
	MediumStressDurationInSeconds      uint
	HighStressDurationInSeconds        uint
	StressQualifier                    string
	StepsGoal                          uint
	NetKilocaloriesGoal                uint
	IntensityDurationGoalInSeconds     uint
	FloorsClimbedGoal                  uint
}

// Dailies mirrors the data structure from an incoming push
type Dailies struct {
	Epochs []Daily
}

// HandleDailies accepts an incoming push from Garmin with Dailies data
func HandleDailies(w http.ResponseWriter, r *http.Request) {
	dailies := Dailies{
		Epochs: make([]Daily, 0),
	}

	err := json.NewDecoder(r.Body).Decode(&dailies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if len(dailies.Epochs) == 0 {
		http.Error(w, "Expected content (epochs) does not contain data", http.StatusBadRequest)
	}
	log.Printf("Incoming Dailies: %v", dailies)
}
