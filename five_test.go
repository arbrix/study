package main

import (
	"testing"
)

func BenchmarkFindEquation(b *testing.B) {
	for n := 0; n < b.N; n++ {
		findEquation([]int{1,2,3,4,5,6,7,8,9}, 100, "")
	}
}

func BenchmarkFindEquation2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		findEquation2(1, 98765432, 100, "")
	}
}

func BenchmarkFindEquation3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		findEquation3(987654321, 100)
	}
}