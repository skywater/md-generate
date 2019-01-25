// CreateMD
package main

import (
	"fmt"
)

func main() {
	path, errors := GetCurrentPath()

	fmt.Println("Hello World!")
	fmt.Println(path)
	fmt.Println(errors)
}
