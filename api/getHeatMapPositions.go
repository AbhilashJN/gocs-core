package api

import (
	"math"
	"os"

	"github.com/AbhilashJN/gocs-core/helpers"
	"github.com/golang/geo/r2"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
	metadata "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/metadata"
)

func GetHeatMapPositions(demoPath string, player string) map[string][]r2.Point {
	//open demo file
	f, err := os.Open(demoPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//init parser
	p := dem.NewParser(f)
	defer p.Close()

	//parse demo headers
	h, herr := p.ParseHeader()
	if herr != nil {
		panic(herr)
	}
	mapName := h.MapName
	mapmeta := metadata.MapNameToMap[mapName]
	reqPlayer := player

	// skip to the end of the knife round
	// knifeRoundEnded := false
	// handlerId := p.RegisterEventHandler(func(e events.RoundEnd) {
	// 	if !knifeRoundEnded {
	// 		knifeRoundEnded = true
	// 	}
	// })
	// for !knifeRoundEnded {
	// 	p.ParseNextFrame()
	// }
	// p.UnregisterEventHandler(handlerId)

	//skip to warmup
	isWarmup := false
	for !isWarmup {
		p.ParseNextFrame()
		isWarmup = p.GameState().IsWarmupPeriod()
	}

	//skip to match
	matchStarted := false
	for !matchStarted {
		p.ParseNextFrame()
		matchStarted = !(p.GameState().IsWarmupPeriod())
	}

	// record deaths
	deathPositions := make(map[string][]r2.Point)
	bombPlantPositions := []r2.Point{}

	deathPositions["kill"] = []r2.Point{}
	deathPositions["death"] = []r2.Point{}
	p.RegisterEventHandler(func(e events.Kill) {
		helpers.DeathTaken(e, reqPlayer, mapmeta, deathPositions)
	})
	p.RegisterEventHandler(func(e events.BombPlanted) {
		bombPlantPositions = helpers.BombPlantPosition(e, reqPlayer, mapmeta, bombPlantPositions)
	})

	// Parse demo to end
	parseErr := p.ParseToEnd()
	if parseErr != nil {
		panic(parseErr)
	}

	normalizedDeathPts := []r2.Point{}
	normalizedKillPts := []r2.Point{}
	normalizedBombPlantPts := []r2.Point{}

	for _, p := range deathPositions["death"] {
		normalizedDeathPts = append(normalizedDeathPts, r2.Point{X: math.Floor(p.X / 2), Y: math.Floor(p.Y / 2)})
	}
	for _, p := range deathPositions["kill"] {
		normalizedKillPts = append(normalizedKillPts, r2.Point{X: math.Floor(p.X / 2), Y: math.Floor(p.Y / 2)})
	}

	for _, p := range bombPlantPositions {
		normalizedBombPlantPts = append(normalizedBombPlantPts, r2.Point{X: math.Floor(p.X / 2), Y: math.Floor(p.Y / 2)})
	}

	finalPositions := make(map[string][]r2.Point)

	finalPositions["kill"] = normalizedKillPts
	finalPositions["death"] = normalizedDeathPts
	finalPositions["bomb_plant"] = normalizedBombPlantPts

	return finalPositions
}
