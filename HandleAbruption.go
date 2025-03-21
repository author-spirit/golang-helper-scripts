package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ch := make(chan os.Signal, 1)

	signal.Notify(ch,
		syscall.SIGINT,  // Any interrupt triggered from keyboard
		syscall.SIGQUIT, // Triggered at ctrl + \ quiting
		syscall.SIGTERM /** Process terminates by kill command*/)

	sig := <-ch
	fmt.Println("\nProgram " + sig.String())
}
