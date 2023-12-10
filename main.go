package main

import (
	"github.com/oakmound/oak/v4"
	"parking-concurrency/scenes"
)

func main() {
	parkingScene := scenes.NewParkingLotScene()

	parkingScene.Start()

	_ = oak.Init("parkingScene", func(c oak.Config) (oak.Config, error) {
		c.BatchLoad = true
		c.Assets.ImagePath = "assets/images"
		c.Title = "Estaciona2"
		return c, nil
	})
}
