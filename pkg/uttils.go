package pkg

func IntInSlice(digit int, digitSlice []int) bool {

	for _, d := range digitSlice {
		if d == digit {
			return true
		}
	}
	return false
}
