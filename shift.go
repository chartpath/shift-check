package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type shift struct {
	Start int64
	End   int64
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
	s.Start, _ = strconv.ParseInt(shiftSourceData.Start, 10, 32)
	s.End, _ = strconv.ParseInt(shiftSourceData.End, 10, 32)
	return nil
}

func main() {
	var userShifts []shift
	errUserShifts := json.Unmarshal(userShiftData, &userShifts)
	if errUserShifts != nil {
		fmt.Println("error:", errUserShifts)
	}
	// fmt.Printf("userShifts, %+v", userShifts)

	var availableShifts []shift
	errAvailShifts := json.Unmarshal(availableShiftData, &availableShifts)
	if errAvailShifts != nil {
		fmt.Println("error:", errAvailShifts)
	}
	// fmt.Printf("availableShifts %+v", availableShifts)

	for _, shift := range availableShifts {
		shiftIsValid := isValidShift(userShifts, shift)
		// fmt.Println("shiftIsValid", shiftIsValid)
		if shiftIsValid {
			fmt.Println("shift is valid", shift)
		} else {
			fmt.Println("shift is invalid", shift)
		}
	}

}

func isValidShift(u []shift, c shift) bool {
	var exitingStarts []int64
	var existingEnds []int64
	for _, shift := range u {
		exitingStarts = append(exitingStarts, shift.Start)
		existingEnds = append(existingEnds, shift.End)

		// todo: probably want to check the date packages for helper functions
	}

	var startOk bool
	var endOk bool

	// candidate start AFTER any user ends AND candidate ends BEFORE all user starts
	for _, exEnd := range existingEnds {
		if c.Start >= exEnd {
			startOk = true
			for _, exStart := range exitingStarts {
				if c.End <= exStart {
					endOk = true
				}
			}
		}
	}

	// OR
	// candidate ends BEFORE any user starts AND candidate starts AFTER all user starts
	for i, exStart := range exitingStarts {
		if c.End <= exStart {
			if i >= 1 && c.End <= exitingStarts[i-1] { // lookbehind 1 deep
				endOk = true
			}
			for _, exEnd := range existingEnds {
				if c.Start >= exEnd || c.Start == 0 {
					startOk = true
				}
			}
		}
	}

	// fmt.Println(startOk, endOk)
	return startOk && endOk
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
