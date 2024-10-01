package main

import (
	"apt-remove-version/cmd"
	"log/slog"
)

func main() {
	if err := cmd.Run(); err != nil {
		slog.Error("", "error", err.Error())
	}
}
