package main

import (
	"fmt"
	"strings"
)

func main() {

	var sc StockCut

	//sc.Demand = Parts{
	//	Part{1100, 20},
	//	Part{750, 60},
	//	Part{2000, 10},
	//}
	//
	//sc.Inventory = Stocks{
	//	Stock{2300, 1},
	//	Stock{3500, 1},
	//	Stock{5000, 3},
	//	Stock{12192, 4},
	//	Stock{15240, 3},
	//	Stock{18288, 3},
	//}

	sc.Demand = Parts{
		Part{1100, 5},
	}

	sc.Inventory = Stocks{
		Stock{2300, 4},
		Stock{3500, 3},
		Stock{6000, 1},
		Stock{1200, 5},
	}

	sc.Kerf = 3

	fmt.Println("Best Fit Decreasing")
	sc.BestFitOptimisation()
	fmt.Println(strings.Repeat("=", 30))

	fmt.Println("Random Best Fit")
	sc.RandomBestFitOptimisation()
	fmt.Println(strings.Repeat("=", 30))

	fmt.Println("First Fit")
	sc.FirstFitOptimisation()
	fmt.Println(strings.Repeat("=", 30))

	/* GA */
	fmt.Println(strings.Repeat("=", 50))

	ga := GAStockCut{
		PopulationSize:  10,
		GenerationToRun: 100,
		Stocks: Stocks{
			Stock{2300, 4},
			Stock{3500, 3},
			Stock{6000, 1},
			Stock{1200, 5},
		},
		Parts: Parts{
			Part{Length: 1100, Quantity: 5},
		},
		Kerf: 3,
	}

	ga.Run()

	//TODO interface to heuristics StockCut or GAStockcut

}
