package main

import (
	"fmt"
	"time"

	"github.com/AurelienS/cigare/internal/webserver"
	"github.com/AurelienS/cigare/pkg/igcparser"
	"github.com/labstack/echo/v4"
)

const igcFile = "/mnt/c/Users/TheGosu/Desktop/9-9-2023--13-47.igc"
const ouputFile = "test.json"
const outputFormat = "json"

func main() {

	start := time.Now()
	flight, _ := igcparser.Parse()
	flight.Initialize()

	flight.Stats.PrettyPrint()
	flight.Draw2DMap(true)
	flight.DrawElevation()
	fmt.Println("file: main.go ~ line 25 ~ elapsed : ", time.Now().Sub(start))

	e := echo.New()
	router := webserver.Router{
		Handler: webserver.Handler{},
	}
	router.Initialize(e)

	e.Logger.Fatal(e.Start(":3000"))

}
