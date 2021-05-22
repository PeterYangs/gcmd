package main

import (
	"fmt"
	"github.com/PeterYangs/cmd"
)

func main() {

	command := cmd.Command("php index.php").WaitCustomChan()

	outChan := command.GetCustomOutChan()
	errChan := command.GetCustomErrChan()

	go func() {

		for t := range outChan {

			fmt.Println(string(t))
		}

		fmt.Println("out chan close")

		command.Done()

	}()

	go func() {

		for t := range errChan {

			fmt.Println("err:", string(t))
		}

		fmt.Println("err chan close")

		command.Done()

	}()

	command.Start()

}
