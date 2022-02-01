package main

import "sort"

type Population []Individual

//Stats return best and average fitness of the population
func (pop Population) Stats() (float64, float64) {
	var best, avg, sum float64

	for i := range pop {
		if pop[i].Fitness > best {
			best = pop[i].Fitness
		}
		sum += pop[i].Fitness
	}

	avg = sum / float64(len(pop))

	return best, avg
}

func (pop *Population) FitnessSort() {
	sort.Sort(pop)
}

//implements sort.Interface

func (pop Population) Len() int {
	return len(pop)
}

func (pop Population) Less(i, j int) bool {
	return pop[i].Fitness < pop[j].Fitness
}

func (pop Population) Swap(i, j int) {
	pop[i], pop[j] = pop[j], pop[i]
}
