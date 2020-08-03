package main

import (
	"math/rand"
	"testing"
	"time"
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

		go dupsHolder(pairChan, dupGroupsChan)

		for _, pair := range shufPairs {
			pairChan <- pair
		}

		close(pairChan)

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

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func benchmarkDupsHolder(i int, b *testing.B) {
	pairs := make([][]string, 0)

	rand.Seed(time.Now().UnixNano())

	for j := 0; j < (i-1)*(i-1); j++ {
		pairs = append(pairs, []string{randSeq(5), randSeq(5)})
	}

	for n := 0; n < b.N; n++ {
		pairChan := make(chan []string)
		dupGroupsChan := make(chan []string, len(pairs))

		go dupsHolder(pairChan, dupGroupsChan)

		for _, pair := range pairs {
			pairChan <- pair
		}
		close(pairChan)

		for group := range dupGroupsChan {
			res := group
			_ = res
		}
	}
}

func BenchmarkDupHolder10(b *testing.B) {
	benchmarkDupsHolder(10, b)
}

func BenchmarkDupHolder100(b *testing.B) {
	benchmarkDupsHolder(100, b)
}
