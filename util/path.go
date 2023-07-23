package util

import "os"

func ShortenPath(path string, levelsLeft uint) string {
	s := string(os.PathSeparator) + "..."
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == os.PathSeparator {
			levelsLeft--
			if levelsLeft == 0 {
				s += path[i:]
				break
			}
		}
	}
	return s
}

func IsPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
