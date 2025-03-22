package main

import (
	"github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd"
	_ "github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd/list"
	_ "github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd/stagerOne"
	_ "github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd/starListen"
)

func main() {
	cmd.Execute()
}
