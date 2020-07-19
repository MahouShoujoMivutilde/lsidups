package main

// ContainsStr tells whether slice contains x.
func ContainsStr(slice []string, x string) bool {
	for _, n := range slice {
		if x == n {
			return true
		}
	}
	return false
}
