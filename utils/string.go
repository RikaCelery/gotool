package utils

func Contains(slice []string, item string) bool {
	for _, value := range slice {
		if value == (item) {
			return true
		}
	}
	return false
}
