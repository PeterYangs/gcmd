package main

import (
	"fmt"
	"github.com/PeterYangs/gcmd"
)

func main() {

	out, err := gcmd.Command("php index.php").OutPut().Start()

	if err != nil {

		fmt.Println(err)

		return
	}

	fmt.Println(string(out))
}
