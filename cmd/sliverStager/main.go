package main

import (
	"github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd"
	_ "github.com/Esonhugh/sliver-stage-helper/cmd/sliverStager/cmd/stagerOne"
)

func main() {
	cmd.Execute()
}
