package main

import (
	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/config"
	"github.com/kelwynOliveira/Goexpert-Rate-Limiter/router"
)

func main() {
	config.Init()
	router.Init()
}
