package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println(os.Args)
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("have not define")
	}
}

func run() {
	cmd := exec.Command(os.Args[2])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}
}
