package api

import (
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
)

func GetMapName(demoPath string) string {
	f, err := os.Open(demoPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//init parser
	p := dem.NewParser(f)
	defer p.Close()

	h, herr := p.ParseHeader()
	if herr != nil {
		panic(herr)
	}
	return h.MapName
}
