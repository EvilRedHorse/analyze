package main

import (
	"gitlab.com/NebulousLabs/analysis/jsontag"
	"gitlab.com/NebulousLabs/analysis/lockcheck"
	"gitlab.com/NebulousLabs/analysis/responsewritercheck"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		lockcheck.Analyzer,
		responsewritercheck.Analyzer,
		jsontag.Analyzer,
	)
}
