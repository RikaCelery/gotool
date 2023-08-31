package utils

import "os"

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func CheckFiles(files []string) (invalids []string) {
	for _, file := range files {
		if !IsExist(file) {
			invalids = append(invalids, file)
		}
	}
	return invalids
}
