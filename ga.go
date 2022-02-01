package main

import (
	"fmt"
)

const EvalWaste = 75.0
const EvalMinStockUsage = 25.0

type GAStockCut struct {
	PopulationSize  int
	GenerationToRun int

	population Population

	Stocks Stocks
	Parts  Parts
	Kerf   uint

	xpPartsLookup []int  //Expanded StockIndex
	xpPartsLength []uint //Expanded StockLength + Kerf

	//elitism() saving best chromosomes in new population
	//recombination()

}

func (ga *GAStockCut) Run() {
	ga.initialize()

	for gen := 1; gen <= ga.GenerationToRun; gen++ {

		ga.evaluatePopulation()
		best, avg := ga.population.Stats()
		fmt.Printf("Generation %03d | Best %.2f | Average %.2f\n", gen, best, avg)
	}
}

func (ga *GAStockCut) initialize() {
	//Expand Parts (one index for each part quantity).
	ga.Parts.SortDecreasing()
	ga.xpPartsLength = make([]uint, ga.Parts.QuantitySum())
	ga.xpPartsLookup = make([]int, len(ga.xpPartsLength))
	k := 0
	for j := 0; j < ga.Parts.Len(); j++ {
		for qty := uint(0); qty < ga.Parts[j].Quantity; qty++ {
			ga.xpPartsLength[k] = ga.Parts[j].Length + ga.Kerf
			ga.xpPartsLookup[k] = j
			k++
		}
	}

	//Expand Stock to Chromosome (one index for each stock quantity).
	chrm := make(Chromosome, ga.Stocks.QuantitySum())
	k = 0
	for j := 0; j < ga.Stocks.Len(); j++ {
		for qty := uint(0); qty < ga.Stocks[j].Quantity; qty++ {
			chrm[k] = j
			k++
		}
	}

	//Create initial population
	if ga.PopulationSize == 0 {
		ga.PopulationSize = 100 //set a default
	}
	ga.population = make(Population, ga.PopulationSize)

	//copy chromosome to individual
	for i := range ga.population {
		ga.population[i].Chrm = make(Chromosome, len(chrm))
		copy(ga.population[i].Chrm, chrm)
		ga.population[i].Chrm.Shuffle() //Shuffle chromosomes to make it unique
	}

	if ga.GenerationToRun == 0 {
		ga.GenerationToRun = 1000 //set a default
	}

}

func (ga *GAStockCut) evaluatePopulation() {
	for i := range ga.population {
		ga.firstFitHeuristics(i)
		//fmt.Println(i, strings.Repeat("-", 25))
		ga.evaluateIndividual(i, false)
	}
}

func (ga *GAStockCut) firstFitHeuristics(ind int) {
	//part index = Inventory used
	ga.population[ind].partStockMapping = make([]int, len(ga.xpPartsLength))
	for i := range ga.population[ind].partStockMapping {
		ga.population[ind].partStockMapping[i] = -1 //-1 mean not placed
	}

	remainingStockLength := make([]uint, len(ga.population[ind].Chrm)) //Available length on expanded stock
	//build available Stock Length from individual chromosomes (Stock index)
	for i, sIdx := range ga.population[ind].Chrm {
		remainingStockLength[i] = ga.Stocks[sIdx].Length
	}

	//FirstFit place the part where it can fit
	for p := range ga.xpPartsLength {
		for s := range remainingStockLength {
			if remainingStockLength[s] >= ga.xpPartsLength[p] {
				ga.population[ind].partStockMapping[p] = s     //link part to inventory index
				remainingStockLength[s] -= ga.xpPartsLength[p] //decrease stock remaining space
				break
			}
		}
	}
}

func (ga *GAStockCut) evaluateIndividual(ind int, verbose bool) {

	StockUsage := make(map[int]uint) //Stock Index vs Total Part length

	//Group stock index + sum usage
	for pIdx, sIdx := range ga.population[ind].partStockMapping {
		//pIdx: expanded Part Position on array
		//sIdx: expanded Stock Position on array

		if sIdx >= 0 {
			StockUsage[sIdx] += ga.Parts[ga.xpPartsLookup[pIdx]].Length //real part length (without kerf)
		} else if verbose {
			fmt.Printf("Part %d not placed\n", pIdx) //TODO add penalty for unplaced total length or count
		}
	}

	//Compute yield
	//Stock usage follow xpStockLookup
	var totalPartsLength, usedStockTotal uint
	for sIdx, usedLength := range StockUsage {
		totalPartsLength += usedLength
		stockLength := ga.Stocks[ga.population[ind].Chrm[sIdx]].Length
		usedStockTotal += stockLength

		//Bar representation
		if verbose {
			fmt.Printf("Bar %d: (L=%d) \t|", ga.population[ind].Chrm[sIdx], stockLength)
			for i := range ga.population[ind].partStockMapping {
				if ga.population[ind].partStockMapping[i] == sIdx {
					fmt.Printf("%d|", ga.Parts[ga.xpPartsLookup[i]].Length) //print part real length
				}
			}
			fmt.Printf("(%d)|\n", stockLength-usedLength) //waste (including kerf)
		}
	}

	wasteScore := float64(totalPartsLength) / float64(usedStockTotal) //Less waste is better
	minStockUsageScore := 1.0 / float64(len(StockUsage))              //Less stock used is the better
	evalTotal := EvalWaste + EvalMinStockUsage

	ga.population[ind].Fitness = EvalWaste/evalTotal*wasteScore + EvalMinStockUsage/evalTotal*minStockUsageScore

	//Global Stats
	if verbose {
		fmt.Printf("Total parts length \t%d\n", totalPartsLength)
		fmt.Printf("Used stocks total \t%d\n", usedStockTotal)
		fmt.Printf("Total yield \t\t%.3f%%\n", float64(totalPartsLength)/float64(usedStockTotal)*100.0)
		fmt.Printf("Score \t\t%.3f%%\n", ga.population[ind].Fitness)
	}

}

//TODO: BestFit Heuristics
