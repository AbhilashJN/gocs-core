package api

import (
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
)

func ListPlayers(demoPath string) []string {
	f, err := os.Open(demoPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//init parser
	p := dem.NewParser(f)
	defer p.Close()

	// // skip to the end of the knife round
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

	allPlayers := p.GameState().Participants().Playing()
	pnames := []string{}

	for _, p := range allPlayers {
		pnames = append(pnames, p.Name)
	}

	return pnames
}
