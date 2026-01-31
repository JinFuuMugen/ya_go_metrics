package main

import (
	"log"
	"os"
)

func main() {
	log.Fatal("x")
	os.Exit(0)
}

func helper() {
	log.Fatal("x") // want "log\\.Fatal is not allowed outside main\\.main"
}

func helper2() {
	os.Exit(1) // want "os\\.Exit is not allowed outside main\\.main"
}

func p() {
	panic("x") // want "avoid using panic"
}
