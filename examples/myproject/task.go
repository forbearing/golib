package main

import "github.com/forbearing/golib/logger"

func SayHello() error {
	// fmt.Println("hello world!")
	logger.Task.Info("hello world!")
	return nil
}

func SayGoodbye() error {
	// fmt.Println("goodbye world!")
	logger.Task.Info("goodbye world!")
	return nil
}
