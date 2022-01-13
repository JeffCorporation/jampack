package main

import (
	"reflect"
	"testing"
)

func TestChromosome_CrossOver(t *testing.T) {
	type args struct {
		parent2    Chromosome
		crossPoint int
	}
	tests := []struct {
		name    string
		parent1 Chromosome
		args    args
		want    Chromosome
	}{
		{
			name:    "Test 1",
			parent1: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			args: args{
				parent2:    []int{2, 4, 6, 8, 0, 1, 3, 5, 7, 9},
				crossPoint: 3,
			},
			want: []int{1, 2, 3, 4, 6, 8, 0, 5, 7, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.parent1.CrossOver(tt.args.parent2, tt.args.crossPoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CrossOver() = %v, want %v", got, tt.want)
			}
		})
	}
}
