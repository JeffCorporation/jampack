package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type StockCut struct {
	Demand    Parts
	Inventory Stocks
	Kerf      uint //Material removed from blad

	xpPartsLength []uint //Expanded StockLength + Kerf
	xpPartsLookup []int  //Expanded StockIndex

}

func (sc *StockCut) initialize() {
	rand.Seed(time.Now().UnixNano())

	//Expand Parts
	sc.Demand.SortDecreasing()
	sc.xpPartsLength = make([]uint, sc.Demand.QuantitySum())
	sc.xpPartsLookup = make([]int, len(sc.xpPartsLength))

	i := 0
	for j := 0; j < sc.Demand.Len(); j++ {
		for qty := uint(0); qty < sc.Demand[j].Quantity; qty++ {
			sc.xpPartsLength[i] = sc.Demand[j].Length + sc.Kerf
			sc.xpPartsLookup[i] = j
			i++
		}
	}

	//Sort Inventory
	sc.Inventory.SortIncreasing()

}

//expandedStockLookup return expanded stock index (one index for each stock quantity).
//Randomize parameter shuffle the list, otherwise it will be sort in ascending order (ordered on initialize function)
func (sc *StockCut) expandedStockLookup(randomize bool) []int {
	lookup := make([]int, sc.Inventory.QuantitySum())

	i := 0
	for j := 0; j < sc.Inventory.Len(); j++ {
		for qty := uint(0); qty < sc.Inventory[j].Quantity; qty++ {
			lookup[i] = j
			i++
		}
	}

	if randomize {
		shuf := rand.Intn(len(lookup))
		for s := 0; s < shuf; s++ {
			a := rand.Intn(len(lookup))
			b := rand.Intn(len(lookup))
			lookup[a], lookup[b] = lookup[b], lookup[a]
		}
	}

	return lookup
}

func (sc *StockCut) BestFitOptimisation() {
	sc.initialize()

	stockIndexLookup := sc.expandedStockLookup(false)
	partStockMapping := sc.bestFit(stockIndexLookup)

	sc.evaluate(partStockMapping, stockIndexLookup, true)
}

func (sc *StockCut) RandomBestFitOptimisation() {
	sc.initialize()

	rndStockIndex := sc.expandedStockLookup(true)
	partStockMapping := sc.bestFit(rndStockIndex)

	sc.evaluate(partStockMapping, rndStockIndex, true)
}

func (sc *StockCut) FirstFitOptimisation() {
	sc.initialize()

	stockIndexLookup := sc.expandedStockLookup(false)
	partStockMapping := sc.firstFit(stockIndexLookup)

	sc.evaluate(partStockMapping, stockIndexLookup, true)
}

//BestFit Use BestFit Decreasing heuristics to place parts on stock using stock order in parameter
//Return a slice indicating the xpStockIndex used for each xpPartIndex
//
//Place th  parts from longest to shortest on the provided Stock sequence.
//Subsequent parts are placed where there is the less waste.
//To use this algorithm, part must be sorted in decreasing order, and stock in ascending.
func (sc *StockCut) bestFit(xpStockIndex []int) []int {
	//start := time.Now()

	//part index = Inventory used
	partsStockMapping := make([]int, len(sc.xpPartsLength))
	for i := range partsStockMapping {
		partsStockMapping[i] = -1 //-1 mean not placed
	}

	remainingStockLength := make([]uint, len(xpStockIndex)) //Available length on expanded stock
	//build available Stock Length from provided xpStockIndex
	for i := range xpStockIndex {
		remainingStockLength[i] = sc.Inventory[xpStockIndex[i]].Length
	}

	//BestFit biggest part in smallest stock
	upperBound := 0 // limit the stock usage
	for p := range sc.xpPartsLength {

		//search for BestFit in currently used stock
		bestFitIdx := -1
		lessWaste := uint(math.MaxUint)

		//BestFit Decreasing
		for s := 0; s <= upperBound; s++ {
			//look for a remaining stock space than can fit the part and generate the less waste
			if (remainingStockLength[s] > sc.xpPartsLength[p]) && ((remainingStockLength[s] - sc.xpPartsLength[p]) < lessWaste) {
				lessWaste = remainingStockLength[s] - sc.xpPartsLength[p]
				bestFitIdx = s //Save Inventory index
				if remainingStockLength[s]-sc.xpPartsLength[p] == 0 {
					break //perfect fit - don't search further
				}
			}

		}

		//Handle Best fit if successful
		if bestFitIdx >= 0 {
			partsStockMapping[p] = bestFitIdx                       //link part to inventory index
			remainingStockLength[bestFitIdx] -= sc.xpPartsLength[p] //decrease stock remaining space
			continue
		}

		//BestFit unsuccessful - Search for a new stock
		for s := upperBound + 1; s < len(remainingStockLength); s++ {

			if sc.xpPartsLength[p] <= remainingStockLength[s] {
				partsStockMapping[p] = s                       //link part to inventory index
				remainingStockLength[s] -= sc.xpPartsLength[p] //decrease stock remaining space
				upperBound = s                                 //new stock limit
				break                                          //break new stock search
			}
		}

	}

	//elapsed := time.Since(start)
	return partsStockMapping
}

func (sc *StockCut) firstFit(xpStockIndex []int) []int {
	//part index = Inventory used
	partsStockMapping := make([]int, len(sc.xpPartsLength))
	for i := range partsStockMapping {
		partsStockMapping[i] = -1 //-1 mean not placed
	}

	remainingStockLength := make([]uint, len(xpStockIndex)) //Available length on expanded stock
	//build available Stock Length from provided xpStockIndex
	for i := range xpStockIndex {
		remainingStockLength[i] = sc.Inventory[xpStockIndex[i]].Length
	}

	//FirstFit place the part where where it can fit
	for p := range sc.xpPartsLength {
		for s := range remainingStockLength {
			if remainingStockLength[s] > sc.xpPartsLength[p] {
				partsStockMapping[p] = s                       //link part to inventory index
				remainingStockLength[s] -= sc.xpPartsLength[p] //decrease stock remaining space
			}
		}
	}

	return partsStockMapping
}

//Evaluate return the yield value and optionally print the results
func (sc *StockCut) evaluate(partStockMapping []int, xpStockLookup []int, verbose bool) float64 {
	StockUsage := make(map[int]uint) //Stock Index vs Total Part length

	//Group stock index + sum usage
	for pIdx, sIdx := range partStockMapping {
		//pIdx: expanded Part Position on array
		//sIdx: expanded Stock Position on array

		if sIdx >= 0 {
			StockUsage[sIdx] += sc.Demand[sc.xpPartsLookup[pIdx]].Length //real part length (without kerf)
		} else if verbose {
			fmt.Printf("Part %d not placed\n", pIdx) //TODO add penalty for unplaced total length or count
		}
	}

	//Compute yield
	//Stock usage follow xpStockLookup
	var totalPartsLength, usedStockTotal uint
	for sIdx, usedLength := range StockUsage {
		totalPartsLength += usedLength
		stockLength := sc.Inventory[xpStockLookup[sIdx]].Length
		usedStockTotal += stockLength

		//Bar representation
		if verbose {
			fmt.Printf("Bar %d: (L=%d) \t|", sIdx, stockLength)
			for i := range partStockMapping {
				if partStockMapping[i] == sIdx {
					fmt.Printf("%d|", sc.Demand[sc.xpPartsLookup[i]].Length) //print part real length
				}
			}
			fmt.Printf("(%d)|\n", stockLength-usedLength) //waste (including kerf)
		}
	}

	//Global Stats
	if verbose {
		fmt.Printf("Total parts length \t%d\n", totalPartsLength)
		fmt.Printf("Used stocks total \t%d\n", usedStockTotal)
		fmt.Printf("Total yield \t\t%.3f%%\n", float64(totalPartsLength)/float64(usedStockTotal)*100.0)
	}

	return float64(totalPartsLength) / float64(usedStockTotal)
}
