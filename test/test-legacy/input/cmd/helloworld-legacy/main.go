package main

import (
	"example.com/foobar"
	helloworld "example.com/helloworld-legacy"
)

func main() {
	println(helloworld.Hello())
	println(foobar.Hello())
}
