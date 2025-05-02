package main

import (
	"github.com/ayayaakasvin/gen-project/internal/app"
	"github.com/ayayaakasvin/gen-project/internal/logger"
)

func main() {
	logger := logger.NewLogger()
	app.App(logger)	
}