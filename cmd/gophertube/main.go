package main

import (
	"context"
	"fmt"
	"gophertube/internal/app"
	"os"
)

func main() {
	gophertube := app.New()
	if err := gophertube.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\n\033[1;33m%v\033[0m\n", err)
		os.Exit(1)
	}
}
