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

// findGroup returns index of the group in dups that contains fp,
// -1 if not found anywhere
func findGroup(dups [][]string, fp string) int {
	for i, group := range dups {
		if ContainsStr(group, fp) {
			return i
		}
	}
	return -1
}
