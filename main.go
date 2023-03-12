package main

import (
	"ginchat/router"
	"ginchat/utils"
)

func main() {
	r := router.Router()
	utils.InitConfig()
	utils.InitMySql()
	utils.InitRedis()
	r.Run()
}
