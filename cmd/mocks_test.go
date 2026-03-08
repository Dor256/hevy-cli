package cmd

import (
	"bytes"
	"context"
	"hevy_cli/internal/hevyapi"
	"github.com/spf13/cobra"
)

type mockHevyAPI struct {
	ListWorkoutsFn  func(ctx context.Context, page, pageSize int) (*hevyapi.ListWorkoutsResponse, error)
	ListRoutinesFn  func(ctx context.Context, page, pageSize int) (*hevyapi.ListRoutinesResponse, error)
	GetRoutineFn    func(ctx context.Context, id string) (*hevyapi.GetRoutineResponse, error)
	UpdateRoutineFn func(ctx context.Context, id string, req *hevyapi.UpdateRoutineRequest) (*hevyapi.UpdateRoutineResponse, error)
}

func (mock *mockHevyAPI) ListWorkouts(ctx context.Context, page, pageSize int) (*hevyapi.ListWorkoutsResponse, error) {
	return mock.ListWorkoutsFn(ctx, page, pageSize)
}

func (mock *mockHevyAPI) ListRoutines(ctx context.Context, page, pageSize int) (*hevyapi.ListRoutinesResponse, error) {
	return mock.ListRoutinesFn(ctx, page, pageSize)
}

func (mock *mockHevyAPI) GetRoutine(ctx context.Context, id string) (*hevyapi.GetRoutineResponse, error) {
	return mock.GetRoutineFn(ctx, id)
}

func (mock *mockHevyAPI) UpdateRoutine(ctx context.Context, id string, req *hevyapi.UpdateRoutineRequest) (*hevyapi.UpdateRoutineResponse, error) {
	return mock.UpdateRoutineFn(ctx, id, req)
}

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}
