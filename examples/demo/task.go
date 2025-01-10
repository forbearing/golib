package main

import "fmt"

func TaskSayHello() error {
	fmt.Println("say hello")
	return nil
}

func TaskSayBye() error {
	fmt.Println("say bye")
	return nil
}
