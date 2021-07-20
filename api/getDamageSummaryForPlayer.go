package api

import (
	"os"

	"github.com/AbhilashJN/gocs-core/helpers"
	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

func GetDamageSummaryForPlayer(demoPath string, player string) []helpers.DamageSummary {
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

	//record damage events
	damageMapByWeaponClass := make(helpers.DamageMap)
	damageMapByWeapon := make(helpers.DamageMap)
	damageMapByPlayer := make(helpers.DamageMap)

	p.RegisterEventHandler(func(e events.PlayerHurt) {
		helpers.DamageEventsByPlayer(e, reqPlayer, damageMapByWeaponClass, damageMapByWeapon, damageMapByPlayer)
	})

	// Parse demo to end
	parseErr := p.ParseToEnd()
	if parseErr != nil {
		panic(parseErr)
	}

	//generate damage summary
	damageSummary := helpers.GenerateDamageSummary(damageMapByWeaponClass, damageMapByWeapon, damageMapByPlayer)
	return damageSummary
}
