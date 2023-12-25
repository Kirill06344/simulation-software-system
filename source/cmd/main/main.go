package main

import (
	"fmt"
	"stewie.com/source/internal"
	"stewie.com/source/internal/sources"
	"stewie.com/source/util"
)

func main() {
	cfg, err := util.LoadConfig("../config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	pd := internal.NewPoissonDistribution(cfg.Lambda)
	sourceStorage := sources.NewSourceStorage(cfg.SourcesAmount, cfg.RequestAmount, pd)
	for {
		request, err := sourceStorage.GenerateRequest()
		if err != nil {
			break
		}
		fmt.Println(request)
	}
}
