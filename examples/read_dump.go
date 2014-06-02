// After running this, execute:
//     md5 res/HelloWorld.class res/Dumped.class
// to ensure that both classes are the same!
package main

import (
	"fmt"
	"os"

	"github.com/jcla1/jclass"
)

func main() {
	f, _ := os.Open("res/HelloWorld.class")
	defer f.Close()

	c, err := class.Parse(f)
	if err != nil {
		panic(err)
	}

	f, _ = os.Create("res/Dumped.class")
	defer f.Close()

	fmt.Println(c.Dump(f))
}
