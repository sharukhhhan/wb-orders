package main

import "wb-l-zero/internal/app"

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
