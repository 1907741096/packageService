package utils

func InArray(val int, arr []int) bool {
	for index := range arr {
		if arr[index] == val {
			return true
		}
	}
	return false
}
