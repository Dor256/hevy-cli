package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"hevy_cli/internal/hevyapi"
	"testing"
)

func TestListRoutines(test *testing.T) {
	test.Run("returns a list of routines", func(test *testing.T) {
		want := []hevyapi.RoutineGet{
			{RoutineBase: hevyapi.RoutineBase{Title: "Push Day"}, ID: "abc-123"},
			{RoutineBase: hevyapi.RoutineBase{Title: "Pull Day"}, ID: "def-456"},
		}
		mock := &mockHevyAPI{
			ListRoutinesFn: func(ctx context.Context, page, pageSize int) (*hevyapi.ListRoutinesResponse, error) {
				return &hevyapi.ListRoutinesResponse{
					Page:      1,
					PageCount: 1,
					Routines:  want,
				}, nil
			},
		}

		cmd := RoutinesCmd(mock)
		output, err := executeCommand(cmd, "list")
		if err != nil {
			test.Fatalf("unexpected error: %v", err)
		}

		var got []hevyapi.RoutineGet
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			test.Fatalf("invalid JSON output: %v", err)
		}
		if len(got) != len(want) {
			test.Fatalf("got %d routines, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i].ID != want[i].ID {
				test.Errorf("routine[%d].ID = %q, want %q", i, got[i].ID, want[i].ID)
			}
			if got[i].Title != want[i].Title {
				test.Errorf("routine[%d].Title = %q, want %q", i, got[i].Title, want[i].Title)
			}
		}
	})

	test.Run("returns an empty list if no routines exist", func(test *testing.T) {
		mock := &mockHevyAPI{
			ListRoutinesFn: func(ctx context.Context, page, pageSize int) (*hevyapi.ListRoutinesResponse, error) {
				return &hevyapi.ListRoutinesResponse{
					Page:      1,
					PageCount: 0,
					Routines:  []hevyapi.RoutineGet{},
				}, nil
			},
		}

		cmd := RoutinesCmd(mock)
		output, err := executeCommand(cmd, "list")
		if err != nil {
			test.Fatalf("unexpected error: %v", err)
		}

		var got []hevyapi.RoutineGet
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			test.Fatalf("invalid JSON output: %v", err)
		}
		if len(got) != 0 {
			test.Errorf("got %d routines, want 0", len(got))
		}
	})

	test.Run("returns error when api fails", func(test *testing.T) {
		mock := &mockHevyAPI{
			ListRoutinesFn: func(ctx context.Context, page, pageSize int) (*hevyapi.ListRoutinesResponse, error) {
				return nil, fmt.Errorf("connection refused")
			},
		}

		cmd := RoutinesCmd(mock)
		_, err := executeCommand(cmd, "list")
		if err == nil {
			test.Fatal("expected error, got nil")
		}
	})
}

func TestGetRoutine(test *testing.T) {
	test.Run("returns error if id is not provided", func(test *testing.T) {
		mock := &mockHevyAPI{}

		cmd := RoutinesCmd(mock)
		_, err := executeCommand(cmd, "get")
		if err == nil {
			test.Fatal("expected error for missing --id, got nil")
		}
	})

	test.Run("returns the routine for a given id", func(test *testing.T) {
		want := hevyapi.RoutineGet{
			RoutineBase: hevyapi.RoutineBase{Title: "Leg Day"},
			ID:          "routine-789",
		}
		mock := &mockHevyAPI{
			GetRoutineFn: func(ctx context.Context, id string) (*hevyapi.GetRoutineResponse, error) {
				if id != "routine-789" {
					test.Errorf("GetRoutine called with id %q, want %q", id, "routine-789")
				}
				return &hevyapi.GetRoutineResponse{Routine: want}, nil
			},
		}

		cmd := RoutinesCmd(mock)
		output, err := executeCommand(cmd, "get", "--id", "routine-789")
		if err != nil {
			test.Fatalf("unexpected error: %v", err)
		}

		var got hevyapi.RoutineGet
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			test.Fatalf("invalid JSON output: %v", err)
		}
		if got.ID != want.ID {
			test.Errorf("ID = %q, want %q", got.ID, want.ID)
		}
		if got.Title != want.Title {
			test.Errorf("Title = %q, want %q", got.Title, want.Title)
		}
	})

	test.Run("returns error if no routine matches the given id", func(test *testing.T) {
		mock := &mockHevyAPI{
			GetRoutineFn: func(ctx context.Context, id string) (*hevyapi.GetRoutineResponse, error) {
				return nil, fmt.Errorf("Hevy API returned status 404")
			},
		}

		cmd := RoutinesCmd(mock)
		_, err := executeCommand(cmd, "get", "--id", "nonexistent")
		if err == nil {
			test.Fatal("expected error, got nil")
		}
	})
}

func TestUpdateRoutine(test *testing.T) {
	test.Run("updates a routine and returns the updated routine", func(test *testing.T) {
		want := hevyapi.RoutineGet{
			RoutineBase: hevyapi.RoutineBase{Title: "Updated Push Day"},
			ID:          "routine-abc",
		}
		mock := &mockHevyAPI{
			UpdateRoutineFn: func(ctx context.Context, id string, req *hevyapi.UpdateRoutineRequest) (*hevyapi.UpdateRoutineResponse, error) {
				if id != "routine-abc" {
					test.Errorf("UpdateRoutine called with id %q, want %q", id, "routine-abc")
				}
				if req.Routine.Title != "Updated Push Day" {
					test.Errorf("UpdateRoutine called with title %q, want %q", req.Routine.Title, "Updated Push Day")
				}
				return &hevyapi.UpdateRoutineResponse{Routine: want}, nil
			},
		}

		body := `{"routine":{"title":"Updated Push Day","notes":"","exercises":[]}}`
		cmd := RoutinesCmd(mock)
		output, err := executeCommand(cmd, "update", "--id", "routine-abc", "--body", body)
		if err != nil {
			test.Fatalf("unexpected error: %v", err)
		}

		var got hevyapi.RoutineGet
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			test.Fatalf("invalid JSON output: %v", err)
		}
		if got.ID != want.ID {
			test.Errorf("ID = %q, want %q", got.ID, want.ID)
		}
		if got.Title != want.Title {
			test.Errorf("Title = %q, want %q", got.Title, want.Title)
		}
	})

	test.Run("returns error on update to non-existent routine id", func(test *testing.T) {
		mock := &mockHevyAPI{
			UpdateRoutineFn: func(ctx context.Context, id string, req *hevyapi.UpdateRoutineRequest) (*hevyapi.UpdateRoutineResponse, error) {
				return nil, fmt.Errorf("Hevy API returned status 404")
			},
		}

		body := `{"routine":{"title":"Whatever","notes":"","exercises":[]}}`
		cmd := RoutinesCmd(mock)
		_, err := executeCommand(cmd, "update", "--id", "nonexistent", "--body", body)
		if err == nil {
			test.Fatal("expected error, got nil")
		}
	})

	test.Run("returns error on update with invalid JSON body", func(test *testing.T) {
		mock := &mockHevyAPI{}

		cmd := RoutinesCmd(mock)
		_, err := executeCommand(cmd, "update", "--id", "routine-abc", "--body", "not-json")
		if err == nil {
			test.Fatal("expected error for invalid JSON, got nil")
		}
	})

	test.Run("returns error when update api call fails", func(test *testing.T) {
		mock := &mockHevyAPI{
			UpdateRoutineFn: func(ctx context.Context, id string, req *hevyapi.UpdateRoutineRequest) (*hevyapi.UpdateRoutineResponse, error) {
				return nil, fmt.Errorf("connection refused")
			},
		}

		body := `{"routine":{"title":"Push Day","notes":"","exercises":[]}}`
		cmd := RoutinesCmd(mock)
		_, err := executeCommand(cmd, "update", "--id", "routine-abc", "--body", body)
		if err == nil {
			test.Fatal("expected error, got nil")
		}
	})
}
