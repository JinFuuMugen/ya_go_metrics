package totest

import (
	"log"
	"os"
)

func F() {
	panic("panic") // want "avoid using panic"

	log.Fatal("x")   // want "log\\.Fatal is not allowed outside main\\.main"
	log.Fatalf("x")  // want "log\\.Fatalf is not allowed outside main\\.main"
	log.Fatalln("x") // want "log\\.Fatalln is not allowed outside main\\.main"

	os.Exit(1) // want "os\\.Exit is not allowed outside main\\.main"
}
