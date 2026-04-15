package main

import (
	"context"
	"os"

	"kumarvv.com/mockdata/configs"
	"kumarvv.com/mockdata/generator"
	"kumarvv.com/mockdata/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.LogErrM("Config file argument is required")
		return
	}

	utils.Log("loading config file %s ...", os.Args[1])
	configFile := os.Args[1]
	config, errs := configs.Load(configFile)
	if errs != nil {
		for _, err := range errs {
			utils.LogErr(err)
		}
		return
	}
	utils.Log("config file loaded successfully")

	if err := generator.Generate(context.Background(), config); err != nil {
		utils.LogErr(err)
		return
	}

	utils.Log("DONE")
	print(config)
}
