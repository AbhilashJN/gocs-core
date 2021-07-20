package api

import (
	"os"

	"github.com/AbhilashJN/gocs-core/helpers"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

func GenerateAccuracySummaryForPlayer(demoPath string, player string) []helpers.AccuracySummary {
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
	weaponFireMap := make(map[string]int)
	weaponShotsHitMap := make(map[string]map[string]int)

	p.RegisterEventHandler(func(e events.WeaponFire) {
		helpers.WeaponFiredByPlayer(e, reqPlayer, weaponFireMap)
	})
	p.RegisterEventHandler(func(e events.PlayerHurt) {
		helpers.WeaponShotsHitByPlayer(e, reqPlayer, weaponShotsHitMap)
	})

	// Parse demo to end
	parseErr := p.ParseToEnd()
	if parseErr != nil {
		panic(parseErr)
	}

	accuracySummary := helpers.GenerateAccuracySummary(weaponFireMap, weaponShotsHitMap)
	return accuracySummary
}
