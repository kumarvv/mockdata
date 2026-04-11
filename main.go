package main

import (
	"fmt"
	"os"

	"kumarvv.com/mockdata/configs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("ERROR: Config file argument is required")
		return
	}

	configFile := os.Args[1]
	config, err := configs.LoadConfig(configFile)
	if err != nil {
		fmt.Print(err)
		return
	}

	print(config)
}
