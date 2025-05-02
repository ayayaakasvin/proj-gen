package main

import (
	"github.com/ayayaakasvin/proj-gen/internal/app"
	"github.com/ayayaakasvin/proj-gen/internal/logger"
)

func main() {
	logger := logger.NewLogger()
	app.App(logger)
}
