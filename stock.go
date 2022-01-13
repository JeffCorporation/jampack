package main

import (
	"math"
	"sort"
)

type Stock struct {
	Length   uint
	Quantity uint
}

type Stocks []Stock

func (s *Stocks) SortIncreasing() {
	sort.Sort(s)
}

func (s Stocks) QuantitySum() uint {
	var sum uint
	for i := range s {
		sum += s[i].Quantity
	}
	return sum
}

func (s Stocks) Longest() uint {
	var maxL uint
	for i := range s {
		if s[i].Length > maxL {
			maxL = s[i].Length
		}
	}
	return maxL
}

func (s Stocks) Shortest() uint {
	minL := uint(math.MaxUint)
	for i := range s {
		if s[i].Length < minL {
			minL = s[i].Length
		}
	}
	return minL
}

//implements sort.Interface

func (s Stocks) Len() int {
	return len(s)
}

func (s Stocks) Less(i, j int) bool {
	return s[i].Length < s[j].Length
}

func (s Stocks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
