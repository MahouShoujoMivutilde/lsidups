package main

import (
	"testing"
)

// https://stackoverflow.com/a/30226442/13291900
func permutations(arr [][]string) [][][]string {
	var helper func([][]string, int)
	res := [][][]string{}

	helper = func(arr [][]string, n int) {
		if n == 1 {
			tmp := make([][]string, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func uniqLen(arr [][]string) int {
	uniq := make(map[string]int)
	for _, g := range arr {
		for _, e := range g {
			uniq[e] += 1
		}
	}
	return len(uniq)
}

func TestDupsHolder(t *testing.T) {
	pairs := [][]string{
		{"aaa", "AAA"},
		{"aaa", "AAAA"},
		{"aaa", "AAAAA"},
		{"AAA", "aaa"},
		{"AAAA", "AAA"},
		{"bbb", "BBB"},
		{"bbb", "BBBBB"},
	}
	trueCount := uniqLen(pairs)
	trueGroups := 2

	for _, shufPairs := range permutations(pairs) {
		pairChan := make(chan []string)
		dupGroupsChan := make(chan []string, len(pairs))
		doneChan := make(chan bool)

		go dupsHolder(pairChan, dupGroupsChan, doneChan)

		for _, pair := range shufPairs {
			pairChan <- pair
		}

		doneChan <- true

		count := 0
		var groups [][]string
		for group := range dupGroupsChan {
			groups = append(groups, group)
			count += len(group)
		}

		if count != trueCount {
			t.Fatalf("dupsHolder has lost some values, want %d, got %d", trueCount, count)
		}

		if len(groups) != trueGroups {
			t.Fatalf("clusterization is incorrect, want %d groups, got %d,"+
				"\n output groups: %v\n input shufPairs: %v",
				trueGroups, len(groups), groups, shufPairs)
		}
	}
}
