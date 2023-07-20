package util

import "fmt"

func ToHumanReadable(number int) string {
	var rounded float32
	var postfix string
	if number > 1_000_000_000 {
		rounded = float32(number) / 1_000_000_000
		postfix = "G"
	} else if number > 1_000_000 {
		rounded = float32(number) / 1_000_000
		postfix = "M"
	} else if number > 1_000 {
		rounded = float32(number) / 1_000_000
		postfix = "K"
	} else {
		return fmt.Sprintf("%d Bytes", number)
	}

	return fmt.Sprintf("%.1f %sB", rounded, postfix)
}
