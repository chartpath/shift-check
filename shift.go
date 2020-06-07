package main

import (
	"encoding/json"
	"fmt"
)

type shift struct {
	Start string
	End   string
}

func main() {
	var userShifts []shift
	// this would be better if the values were ints right off the bat
	err := json.Unmarshal(userShiftData, &userShifts)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Printf("userShifts, %+v", userShifts)

	var availableShifts []shift
	// how to not reuse var names for error handling?
	e := json.Unmarshal(availableShiftData, &availableShifts)
	if e != nil {
		fmt.Println("error:", e)
	}
	// fmt.Printf("availableShifts %+v", availableShifts)

	// validShifts := make([]shift, len(availableShifts))
	for _, shift := range availableShifts {
		shiftIsValid := isValidShift(userShifts, shift)
		// fmt.Println("shiftIsValid", shiftIsValid)
		if shiftIsValid {
			fmt.Println("shift is valid", shift)
			// append(validShifts, shift)
		} else {
			fmt.Println("shift is invalid", shift)
		}
	}

}

func isValidShift(u []shift, c shift) bool {
	var exitingStarts []string
	var existingEnds []string
	for _, shift := range u {
		exitingStarts = append(exitingStarts, shift.Start)
		existingEnds = append(existingEnds, shift.End)

		// lexical comparison seems fine, but we could do real math instead
		// probably want to check the date packages

		// uStart, _ := strconv.ParseInt(shift.Start, 10, 32)
		// uEnd, _ := strconv.ParseInt(shift.End, 10, 32)
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
				if c.Start >= exEnd || c.Start == "0000" {
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
	{"start": "1200", "end": "1800"}
]`)
