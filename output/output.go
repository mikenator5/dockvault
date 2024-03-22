package output

import (
	"fmt"
)

// PrintProgressBar prints a progress bar with the given percentage.
// The progress bar is displayed as a string of '=' characters representing the completed portion,
// followed by spaces representing the remaining portion.
// The percentage is displayed next to the progress bar.
func PrintProgressBar(percentage int) {
	barLength := 50
	numBars := int(float64(barLength) * (float64(percentage) / 100))
	bar := "[" + RepeatStr("=", numBars) + RepeatStr(" ", barLength-numBars) + "]"
	fmt.Printf("\r%s %d%%", bar, percentage)
}

// RepeatStr repeats the given string `s` `count` number of times and returns the concatenated result.
func RepeatStr(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
