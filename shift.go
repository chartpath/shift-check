package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

const timeLayout = "1504"

type shift struct {
	Start time.Time
	End   time.Time
}

type shiftSource struct {
	Start string // `json:"start"`
	End   string // `json:"end"`
}

func (s *shift) UnmarshalJSON(j []byte) error {
	var shiftSourceData shiftSource
	err := json.Unmarshal(j, &shiftSourceData)
	if err != nil {
		return err
	}
	// fmt.Println(shiftSourceData)
	s.Start, _ = time.Parse(timeLayout, shiftSourceData.Start)
	s.End, _ = time.Parse(timeLayout, shiftSourceData.End)
	return nil
}

func main() {
	var userShifts []shift
	errUserShifts := json.Unmarshal(userShiftData, &userShifts)
	if errUserShifts != nil {
		fmt.Println("error:", errUserShifts)
	}
	sort.Slice(userShifts, func(i, j int) bool {
		return userShifts[i].Start.Before(userShifts[j].Start)
	})
	// for _, shift := range userShifts {
	// 	fmt.Println(
	// 		"userShift start", shift.Start.Format(timeLayout),
	// 		"end", shift.End.Format(timeLayout))
	// }
	// fmt.Printf("userShifts, %+v", userShifts)

	midnight, _ := time.Parse(timeLayout, "0000")
	var unscheduledPeriods []shift
	for i, uShift := range userShifts {
		if i == 0 && uShift.Start.After(midnight) {
			firstPeriod := shift{midnight, uShift.Start}
			unscheduledPeriods = append(unscheduledPeriods, firstPeriod)
		} else {
			nextPeriod := shift{userShifts[i-1].End, uShift.Start}
			unscheduledPeriods = append(unscheduledPeriods, nextPeriod)
		}
		if i == len(userShifts)-1 {
			lastPeriod := shift{uShift.End, midnight.Add(-time.Second)}
			unscheduledPeriods = append(unscheduledPeriods, lastPeriod)
		}
	}
	// for _, period := range unscheduledPeriods {
	// 	fmt.Println(
	// 		"unscheduledPeriod start", period.Start.Format(timeLayout),
	// 		"end", period.End.Format(timeLayout))
	// }

	var availableShifts []shift
	errAvailShifts := json.Unmarshal(availableShiftData, &availableShifts)
	if errAvailShifts != nil {
		fmt.Println("error:", errAvailShifts)
	}
	sort.Slice(availableShifts, func(i, j int) bool {
		return availableShifts[i].Start.Before(availableShifts[j].Start)
	})
	// for _, shift := range availableShifts {
	// 	fmt.Println(
	// 		"availableShift start", shift.Start.Format(timeLayout),
	// 		"end", shift.End.Format(timeLayout))
	// }
	// fmt.Printf("availableShifts %+v", availableShifts)

	for _, shift := range availableShifts {
		shiftIsValid := isValidShift(unscheduledPeriods, shift)
		// fmt.Println("shiftIsValid", shiftIsValid)
		if shiftIsValid {
			fmt.Println(
				"VALID shift start", shift.Start.Format(timeLayout),
				"end", shift.End.Format(timeLayout))
		} else {
			// fmt.Println(
			// 	"INVALID shift start", shift.Start.Format(timeLayout),
			// 	"end", shift.End.Format(timeLayout))
		}
	}

}

func isValidShift(unscheduledPeriods []shift, candidate shift) bool {
	for _, period := range unscheduledPeriods {
		if candidate.Start.After(period.Start.Add(-time.Minute)) &&
			candidate.End.Before(period.End.Add(time.Minute)) {
			return true
		}
	}
	return false
}

var userShiftData = []byte(`[
	{"start": "0600", "end": "1000"},
	{"start": "1600", "end": "2000"}
]`)

var availableShiftData = []byte(`[
	{"start": "0000", "end": "0500"},
	{"start": "1000", "end": "1600"},
	{"start": "0000", "end": "2359"},
	{"start": "0600", "end": "1800"},
	{"start": "0000", "end": "1200"},
	{"start": "0600", "end": "1200"},
	{"start": "1800", "end": "2359"},
	{"start": "0000", "end": "0600"},
	{"start": "1200", "end": "2359"},
	{"start": "1200", "end": "1800"},
	{"start": "1200", "end": "1800"}
]`)
