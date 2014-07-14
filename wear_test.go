package main

import "fmt"

func ExampleClothes() {
	c := clothes(65)
	d := clothes(80)
	fmt.Println(c)
	fmt.Println(d)
	//output: [shirt trousers]
	// [t-shirt trousers]
}
