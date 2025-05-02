package prompt

import (
	"strings"

	"github.com/manifoldco/promptui"
)

func AskForOverwrite() (bool) {
	prompt := promptui.Prompt{
		Label:     "Project already exists. Do you want to overwrite it?",
		IsConfirm: true,
		Default: "n",
	}

	result, err := prompt.Run()
	if err != nil {
		return false
	}

	return strings.ToLower(result) == "y"
}