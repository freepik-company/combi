package actions

import (
	"os"
	"os/exec"

	"gcmerge/internal/config"
)

func RunActions(actions *[]config.ActionT) (err error) {
	for _, action := range *actions {
		command := exec.Command(action.Command[0], action.Command[1:]...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err = command.Run(); err != nil {
			return err
		}
	}
	return err
}
