package main

import (
	"gitlab.com/NebulousLabs/analyze/jsontag"
	"gitlab.com/NebulousLabs/analyze/lockcheck"
	"gitlab.com/NebulousLabs/analyze/responsewritercheck"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		lockcheck.Analyzer,
		responsewritercheck.Analyzer,
		jsontag.Analyzer,
	)
}
