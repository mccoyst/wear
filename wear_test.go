package main

import "fmt"

func ExampleClothes() {
	cs := clothes(65)
	ds := clothes(80)
	for _, c := range cs {
		fmt.Println(c)
	}
	for _, d := range ds {
		fmt.Println(d)
	}
	//output: [shirt trousers]
	// [t-shirt trousers]
}
