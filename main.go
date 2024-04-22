package main

import (
	"fmt"
)

func main() {
/* Input
 var numero int
 fmt.Println("Por favor ingresa un numero: ")
 fmt.Scanln(&numero)
 if numero == 3{
	fmt.Println("Es el mismo numero el cual estaba pensando")
 }else{
	fmt.Println("Nice try")
 }
 */
 var arr = [...]int{1,2,3,4,5,6,7,8,9}
 for _, elemento := range arr{
	fmt.Println("El Elemento es" , elemento)
 }
}
