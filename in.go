package main

// containsStr tells whether slice contains x.
func containsStr(slice []string, x string) bool {
	for _, n := range slice {
		if x == n {
			return true
		}
	}
	return false
}
