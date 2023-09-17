package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jixunmoe-go/qrc/internal/qrc"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		println("Usage: qrc-decode <input> <output>")
		os.Exit(1)
	}

	r, err := os.Open(args[0])
	if err != nil {
		fmt.Printf("failed to open input: %v", err)
		os.Exit(2)
	}
	defer r.Close()
	w, err := os.Create(args[1])
	if err != nil {
		fmt.Printf("failed to open output: %v", err)
		os.Exit(3)
	}
	defer w.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		fmt.Printf("failed to read input: %v", err)
		os.Exit(4)
	}
	decrypted, err := qrc.DecodeQRC(data)
	if err != nil {
		fmt.Printf("failed to decode: %v", err)
		os.Exit(5)
	}
	_, err = w.Write(decrypted)
	if err != nil {
		fmt.Printf("failed to write output: %v", err)
		os.Exit(6)
	}

	fmt.Printf("done\n")
	os.Exit(0)
}
