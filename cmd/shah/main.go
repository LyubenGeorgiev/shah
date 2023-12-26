package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/LyubenGeorgiev/shah/application"
)

func main() {
	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		fmt.Println("Failed to start app:", err)
	}
}
