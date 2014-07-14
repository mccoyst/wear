package main

import (
	"fmt"
	"testing"
)

func ExampleClothes() {
	cs := clothes(65)
	for _, c := range cs {
		fmt.Println(c)
	}
	//output: [shirt trousers]
}

func TestPossibilities(t *testing.T) {
	c := possibilities(100)
	if len(c) != 2 {
		t.Fatalf("wrong: %v\n", c)
	}

	
}
