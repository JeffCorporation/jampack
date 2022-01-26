package main

import "math/rand"

type Chromosome []int

//CrossOver Perform Uniform cross over.
//The first gene define by a cross point value will be preserve.
//The rest of the chromosome will be fill the second chromosome (parent 2), making sure not to have duplicates genes.
//If parents are different size, the result will be inconsistent.
func (parent1 Chromosome) CrossOver(parent2 Chromosome, crossPoint int) Chromosome {
	unique := make(map[int]struct{}) //Track already place gene (struct take no place)

	//Make sure we dont "index out of range" if parent is different size
	minLen := len(parent1)
	if len(parent2) < minLen {
		minLen = len(parent2)
	}

	//Make sure cross point is valid
	if crossPoint > minLen || crossPoint < 0 {
		crossPoint = 0 // Invalid Cross point - Return parent 1
	}

	//The result
	offsping := make(Chromosome, minLen)

	//Keep the first genes of parent 1
	for i := 0; i < crossPoint; i++ {
		offsping[i] = parent1[i]
		unique[parent1[i]] = struct{}{}
	}

	//Fill the rest of the chromosome with parent 2
	p := crossPoint //track offspring gene position
	for i := range parent2 {
		if _, exist := unique[parent2[i]]; !exist {
			offsping[p] = parent2[i]
			unique[parent2[i]] = struct{}{}
			p++
		}
	}

	return offsping
}

//Shuffle the chromosome order. For initial population.
func (c *Chromosome) Shuffle() {
	for i := 0; i < len(*c); i++ {
		a := rand.Intn(len(*c))
		b := rand.Intn(len(*c))
		(*c)[a], (*c)[b] = (*c)[b], (*c)[a]
	}
}
