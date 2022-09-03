package main

import (
	"miner/config"
	"miner/db"
	"miner/web"
)

func main() {
	config.Init()
	db.Init()
	web.Init()
}
