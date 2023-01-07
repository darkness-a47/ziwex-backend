package utils

func ListContainsString(l []string, s string) bool {
	for _, v := range l {
		if s == v {
			return true
		}
	}
	return false
}
