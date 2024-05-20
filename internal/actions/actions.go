package actions

import (
	"os"
	"os/exec"

	"combi/api/v1alpha1"
)

func RunActions(actions *[]v1alpha1.ActionT) (err error) {
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
