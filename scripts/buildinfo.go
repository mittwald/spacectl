package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	out, err := os.Create("buildinfo/buildinfo.go")
	if err != nil {
		panic(err)
	}

	defer out.Close()

	tag := os.Getenv("CI_BUILD_TAG")
	if tag == "" {
		tag = "nightly"
	}

	hash := os.Getenv("CI_COMMIT_SHA")
	if hash == "" {
		hash = "<local>"
	}

	out.WriteString("package buildinfo\n\nconst (\n")
	fmt.Fprintf(out, "\tVersion = `%s`\n", tag)
	fmt.Fprintf(out, "\tHash = `%s`\n", hash)
	fmt.Fprintf(out, "\tBuildDate = `%s`\n", time.Now().Format(time.RFC1123))
	out.WriteString(")\n")
}
