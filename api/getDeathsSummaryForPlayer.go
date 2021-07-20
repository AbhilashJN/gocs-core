package api

import (
	"os"

	"github.com/AbhilashJN/gocs-core/helpers"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

func GetDeathsSummaryForPlayer(demoPath string, player string) []helpers.DeathsSummary {
	//open demo file
	f, err := os.Open(demoPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//init parser
	p := dem.NewParser(f)
	defer p.Close()

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
	deathsMapByPlayer := make(map[string]map[string]int)
	deathsMapByWeapon := make(map[string]map[string]int)
	deathsMapByWeaponClass := make(map[string]map[string]int)
	p.RegisterEventHandler(func(e events.Kill) {
		helpers.DeathEvents(e, reqPlayer, deathsMapByPlayer, deathsMapByWeapon, deathsMapByWeaponClass)
	})

	// Parse demo to end
	parseErr := p.ParseToEnd()
	if parseErr != nil {
		panic(parseErr)
	}

	deathsSummary := helpers.GenerateDeathsSummary(deathsMapByPlayer, deathsMapByWeapon, deathsMapByWeaponClass)
	return deathsSummary
}
