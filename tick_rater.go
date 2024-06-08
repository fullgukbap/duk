package duk

import (
	"log"
	"time"
)

type TickRater struct {
	app *App

	tickRate     int
	tickRateData any
}

func newTickRater(app *App) *TickRater {
	return &TickRater{
		app: app,

		tickRate:     10,
		tickRateData: struct{}{},
	}
}

func (app *App) Update() {
	interval := time.Minute / time.Duration(app.tickRater.tickRate)
	ticker := time.Tick(interval)

	for range ticker {
		err := app.orchestration.broadcast(app.tickRater.tickRateData)
		if err != nil {
			log.Println("boradcast 실패")
		}
	}
}

func (app *App) Register(tickRate int, tickRateData any) {
	app.tickRater.tickRate = tickRate
	app.tickRater.tickRateData = tickRateData
}
