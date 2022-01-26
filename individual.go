package main

type Individual struct {
	Fitness          float64
	Chrm             Chromosome //Represent a solution - a sequence of stock to be applied to the immutable part sequence
	partStockMapping []int
	//BestFitHeuristics
	//FirstFitHeuristics
	//fitness()
}
