package output

import "fmt"

func PrintUsage() {
	fmt.Println("Usage: dockerHelpers <upload | list | load>")
}

func PrintProgressBar(percentage int) {
	barLength := 50
	numBars := int(float64(barLength) * (float64(percentage) / 100))
	bar := "[" + RepeatStr("=", numBars) + RepeatStr(" ", barLength-numBars) + "]"
	fmt.Printf("\r%s %d%%", bar, percentage)
}

func RepeatStr(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
