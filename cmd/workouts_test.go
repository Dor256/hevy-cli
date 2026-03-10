package cmd

import (
	"context"
	"encoding/json"
	"hevy_cli/internal/hevyapi"
	"testing"
	"time"
)

func TestListWorkouts(test *testing.T) {
	test.Run("returns a list of workouts for a given user", func(test *testing.T) {
		want := []hevyapi.Workout{
			{
				ID:          "abc-123",
				Title:       "workout",
				RoutineID:   "routine-123",
				Description: "some descripiton",
				StartTime:   time.Now(),
				EndTime:     time.Now().Add(time.Duration(1) * time.Hour),
				UpdatedAt:   time.Now(),
				CreatedAt:   time.Now(),
				Exercises: []hevyapi.WorkoutExercise{
					{ExerciseBase: hevyapi.ExerciseBase{Title: "Bench Press"}},
					{ExerciseBase: hevyapi.ExerciseBase{Title: "Squat"}},
				},
			},
		}
		mock := &mockHevyAPI{
			ListWorkoutsFn: func(ctx context.Context, page, pageSize int) (*hevyapi.ListWorkoutsResponse, error) {
				return &hevyapi.ListWorkoutsResponse{
					Page:      1,
					PageCount: 1,
					Workouts:  want,
				}, nil
			},
		}
		cmd := WorkoutsCmd(mock)
		output, err := executeCommand(cmd, "list")
		if err != nil {
			test.Fatalf("unexpected error: %v", err)
		}
		var got []hevyapi.Workout
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			test.Fatalf("invalid JSON output: %v", err)
		}
		if len(got) != len(want) {
			test.Fatalf("got %d workouts, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i].ID != want[i].ID {
				test.Errorf("workout[%d].ID = %q, want %q", i, got[i].ID, want[i].ID)
			}
			if got[i].Title != want[i].Title {
				test.Errorf("workout[%d].Title = %q, want %q", i, got[i].Title, want[i].Title)
			}
		}
	})
	test.Run("returns an empty list if no workouts exist", func(test *testing.T) {
		mock := &mockHevyAPI{
			ListWorkoutsFn: func(ctx context.Context, page, pageSize int) (*hevyapi.ListWorkoutsResponse, error) {
				return &hevyapi.ListWorkoutsResponse{
					Page:      1,
					PageCount: 0,
					Workouts:  []hevyapi.Workout{},
				}, nil
			},
		}

		cmd := WorkoutsCmd(mock)
		output, err := executeCommand(cmd, "list")
		if err != nil {
			test.Fatalf("unexpected error: %v", err)
		}

		var got []hevyapi.Workout
		if err := json.Unmarshal([]byte(output), &got); err != nil {
			test.Fatalf("invalid JSON output: %v", err)
		}
		if len(got) != 0 {
			test.Errorf("got %d workouts, want 0", len(got))
		}
	})
}
