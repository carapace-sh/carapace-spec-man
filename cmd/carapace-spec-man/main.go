package main

import "github.com/rsteube/carapace-spec-man/cmd/carapace-spec-man/cmd"

var version = "develop"

func main() {
	cmd.Execute(version)
}
