package main

import (
	"math"
	"sort"
)

type Part struct {
	Length   uint
	Quantity uint
}

type Parts []Part

func (p *Parts) SortDecreasing() {
	sort.Sort(sort.Reverse(p))
}

func (p Parts) QuantitySum() uint {
	var sum uint
	for i := range p {
		sum += p[i].Quantity
	}
	return sum
}

func (p Parts) Longest() uint {
	var maxL uint
	for i := range p {
		if p[i].Length > maxL {
			maxL = p[i].Length
		}
	}
	return maxL
}

func (p Parts) Shortest() uint {
	minL := uint(math.MaxUint)
	for i := range p {
		if p[i].Length < minL {
			minL = p[i].Length
		}
	}
	return minL
}

//implements sort.Interface

func (p Parts) Len() int {
	return len(p)
}

func (p Parts) Less(i, j int) bool {
	return p[i].Length < p[j].Length
}

func (p Parts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
