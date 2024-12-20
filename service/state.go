package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/vdub-go/model"
)

// Get or initialize state if not exists
func GetState(ctx context.Context, taskDir string, initialState model.TaskState) (model.TaskState, error) {
	statePath := genStatePath(taskDir)

	state := model.TaskState{}

	stateJson, err := os.ReadFile(statePath)
	if err == nil {
		json.Unmarshal(stateJson, &state)
		return state, nil
	}

	cmd := exec.Command("mkdir", "-p", taskDir)
	_, err = cmd.Output()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return state, err
	}

	initialState.Status = "initialized"
	SaveState(ctx, taskDir, initialState)

	return initialState, nil
}

func SaveState(ctx context.Context, taskDir string, state model.TaskState) error {
	statePath := genStatePath(taskDir)

	stateJson, _ := json.MarshalIndent(state, " ", "  ")

	err := os.WriteFile(statePath, stateJson, 0644)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}

func SaveStateStatus(ctx context.Context, taskDir string, state *model.TaskState, newStatus string) error {
	state.Status = newStatus
	return SaveState(ctx, taskDir, *state)
}

func genStatePath(taskDir string) string {
	return fmt.Sprintf("%s/state.json", taskDir)
}
