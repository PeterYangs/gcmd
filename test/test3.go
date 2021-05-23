package main

import (
	"github.com/PeterYangs/gcmd"
)

func main() {

	gcmd.Command("php index.php").ThrowPanic().Start()

}
