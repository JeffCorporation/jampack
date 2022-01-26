package main

type GAStockCut struct {
	PopulationSize  int
	GenerationCount int

	population Population
	generation int

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
		copy(ga.population[i].Chrm, chrm)
		ga.population[i].Chrm.Shuffle()
	}

	if ga.GenerationCount == 0 {
		ga.GenerationCount = 1000 //set a default
	}

}

func (ga *GAStockCut) firstFitHeuristics(ind *Individual) {
	ind.Fitness = 0

	//part index = Inventory used
	ind.partStockMapping = make([]int, len(ga.xpPartsLength))
	for i := range ind.partStockMapping {
		ind.partStockMapping[i] = -1 //-1 mean not placed
	}

	remainingStockLength := make([]uint, len(ind.Chrm)) //Available length on expanded stock
	//build available Stock Length from individual chromosomes (Stock index)
	for _, i := range ind.Chrm {
		remainingStockLength[i] = ga.Stocks[i].Length
	}

	//FirstFit place the part where it can fit
	for p := range ga.xpPartsLength {
		for s := range remainingStockLength {
			if remainingStockLength[s] > ga.xpPartsLength[p] {
				ind.partStockMapping[p] = s                    //link part to inventory index
				remainingStockLength[s] -= ga.xpPartsLength[p] //decrease stock remaining space
			}
		}
	}

}
