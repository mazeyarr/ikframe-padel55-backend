package util

func RemoveElementByIndex[T any](slice []T, i int) []T {
	return append(slice[:i], slice[i+1:]...)
}
