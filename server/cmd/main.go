package main

import (
	"log"
	"realTime/config"
	"realTime/server/api"
)

func main() {
	Env := config.Env
	Env = config.InitConfig("main")
	log.Fatal(api.NewApi(Env.Port).Run())
}
