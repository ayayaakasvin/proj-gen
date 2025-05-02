package app

import (
	"flag"
	"os"

	"github.com/ayayaakasvin/gen-project/internal/generate"
	"github.com/ayayaakasvin/gen-project/internal/options"
	"github.com/sirupsen/logrus"
)

func App (logger *logrus.Logger) {
	opts, afterParse := options.ParseFlags()

	logger.Info("Parsing flags")
	flag.Parse()
	afterParse(opts)
	logger.Info(opts.String())


	logger.Info("Generating project")
	generate.GenerateProject(opts, logger)
	logger.Info("Project generated successfully")
	os.Exit(0)
}