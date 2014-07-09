package main

import "time"

func main() {
	var channel = make(chan string)

	println("1")

	go func() {
		message := <-channel
		println("thread 1 " + message)
	}()

	println("2")

	go func() {
		message := <-channel
		println("thread 2 " + message)
	}()

	println("3")

	channel <- "a"
	channel <- "b"

	time.Sleep(10 * time.Second)
}

