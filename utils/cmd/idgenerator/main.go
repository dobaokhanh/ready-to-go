package main

import (
	"fmt"
	idgenerator "utils/id_generator"
	lockfree "utils/idgeneratorlockfree"
)

func main() {
	var i int = 0
	for i < 100 {
		lockID, _ := idgenerator.NextID()
		fmt.Printf("Lock ID %d \n", lockID)
		lockFree, _ := lockfree.NextID()
		fmt.Printf("Lock free ID %d \n", lockFree)
		i++
	}
}
