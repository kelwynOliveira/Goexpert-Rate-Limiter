package main

import (
	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/config"
	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/pkg/dependencyinjector"
)

func main() {
	configs, err := config.Load(".")
	if err != nil {
		panic(err)
	}

	di := dependencyinjector.NewDependencyInjector(configs)

	deps, err := di.Inject()
	if err != nil {
		panic(err)
	}

	deps.WebServer.Start()
}
