package main

import (
	"crawler/cmd"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 5 {
		log.Println("Please enter options")
		os.Exit(1)
	}
	novelName := os.Args[1]
	fileName := os.Args[2]
	total, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Println("Final chapter required")
		os.Exit(1)
	}
	batchSize, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Println("error max")
		os.Exit(1)
	}
	cmd.Execute(fileName, novelName, total, batchSize)
}
