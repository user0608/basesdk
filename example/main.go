package main

import (
	"basesdk/setup"
	"basesdk/setup/migrations"
	"context"
	"fmt"
	"log/slog"

	"go.uber.org/fx"
)

func main() {
	service := setup.NewService()
	service.Run(fx.Invoke(func(runner *migrations.MigrationRunner) {
		fmt.Println("OK")
		if err := runner.Status(context.Background()); err != nil {
			slog.Error(err.Error())
		}
	}))
}
