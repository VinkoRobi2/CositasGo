package main

import (
	"fmt"
)

func main() {
 var arr  = [...]float64{1,2,3,4,5,6,7,8,9}
 for _,elemento := range arr{
	sum := elemento /2
	fmt.Println(sum)
 }
}

