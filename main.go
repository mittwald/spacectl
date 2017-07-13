package main

import "github.com/mittwald/spacectl/cmd"

//go:generate go run scripts/buildinfo.go

func main() {
	cmd.Execute()
}
