package hevyapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hevy_cli/internal/middleware"
	"io"
	"net/http"
	"time"
)

type HevyAPI interface {
	ListWorkouts(ctx context.Context, page, pageSize int) (*ListWorkoutsResponse, error)
	ListRoutines(ctx context.Context, page, pageSize int) (*ListRoutinesResponse, error)
	GetRoutine(ctx context.Context, id string) (*GetRoutineResponse, error)
	UpdateRoutine(ctx context.Context, id string, req *UpdateRoutineRequest) (*UpdateRoutineResponse, error)
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		baseURL: "https://api.hevyapp.com",
		httpClient: &http.Client{
			Transport: &middleware.AuthTransport{
				Base: http.DefaultTransport,
			},
		},
	}
}

func doJSON[R any](ctx context.Context, httpClient *http.Client, method, url string, body any) (*R, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("Error serializing request body.")
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("Creating request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Executing request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Hevy API returned status %d", resp.StatusCode)
	}
	var result R
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Decoding response: %w", err)
	}
	return &result, nil
}

type RepRange struct {
	Start *int `json:"start"`
	End   *int `json:"end"`
}

type SetBase struct {
	Type            string    `json:"type"`
	WeightKg        *float64  `json:"weight_kg"`
	Reps            *int      `json:"reps"`
	RepRange        *RepRange `json:"rep_range"`
	DistanceMeters  *float64  `json:"distance_meters"`
	DurationSeconds *float64  `json:"duration_seconds"`
	CustomMetric    *float64  `json:"custom_metric"`
}

type RoutineSetGet struct {
	SetBase
	Index int      `json:"index"`
	RPE   *float64 `json:"rpe"`
}

type ExerciseBase struct {
	Title              string `json:"title"`
	ExerciseTemplateID string `json:"exercise_template_id"`
	SupersetsId        *int   `json:"supersets_id"`
	RestSeconds        *int   `json:"rest_seconds"`
	Notes              string `json:"notes"`
}

type RoutineExerciseGet struct {
	ExerciseBase
	Index int             `json:"index"`
	Sets  []RoutineSetGet `json:"sets"`
}

type RoutineBase struct {
	Title string `json:"title"`
}

type RoutineGet struct {
	RoutineBase
	ID        string               `json:"id"`
	FolderID  *int                 `json:"folder_id"`
	UpdatedAt string               `json:"updated_at"`
	CreatedAt string               `json:"created_at"`
	Exercises []RoutineExerciseGet `json:"exercises"`
}

type ListRoutinesResponse struct {
	Page      int          `json:"page"`
	PageCount int          `json:"page_count"`
	Routines  []RoutineGet `json:"routines"`
}

func (client *Client) ListRoutines(ctx context.Context, page, pageSize int) (*ListRoutinesResponse, error) {
	url := fmt.Sprintf("%s/v1/routines?page=%d&pageSize=%d", client.baseURL, page, pageSize)
	return doJSON[ListRoutinesResponse](ctx, client.httpClient, http.MethodGet, url, nil)
}

type GetRoutineResponse struct {
	Routine RoutineGet `json:"routine"`
}

func (client *Client) GetRoutine(ctx context.Context, id string) (*GetRoutineResponse, error) {
	url := fmt.Sprintf("%s/v1/routines/%s", client.baseURL, id)
	return doJSON[GetRoutineResponse](ctx, client.httpClient, http.MethodGet, url, nil)
}

type RoutineExerciseUpdate struct {
	ExerciseBase
	Sets []SetBase `json:"sets"`
}

type RoutineUpdate struct {
	RoutineBase
	Notes     string                  `json:"notes"`
	Exercises []RoutineExerciseUpdate `json:"exercises"`
}

type UpdateRoutineRequest struct {
	Routine RoutineUpdate `json:"routine"`
}

type UpdateRoutineResponse struct {
	Routine RoutineGet `json:"routine"`
}

func (client *Client) UpdateRoutine(ctx context.Context, id string, request *UpdateRoutineRequest) (*UpdateRoutineResponse, error) {
	url := fmt.Sprintf("%s/v1/routines/%s", client.baseURL, id)
	return doJSON[UpdateRoutineResponse](ctx, client.httpClient, http.MethodPut, url, request.Routine)
}

type WorkoutSet struct {
	SetBase
	Index int      `json:"index"`
	RPE   *float64 `json:"rpe"`
}

type WorkoutExercise struct {
	ExerciseBase
	Index int          `json:"index"`
	Sets  []WorkoutSet `json:"sets"`
}

type Workout struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	RoutineID   string            `json:"routine_id"`
	Description string            `json:"description"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CreatedAt   time.Time         `json:"created_at"`
	Exercises   []WorkoutExercise `json:"exercises"`
}

type ListWorkoutsResponse struct {
	Page      int       `json:"page"`
	PageCount int       `json:"page_count"`
	Workouts  []Workout `json:"workouts"`
}

func (client *Client) ListWorkouts(ctx context.Context, page, pageSize int) (*ListWorkoutsResponse, error) {
	url := fmt.Sprintf("%s/v1/workouts?page=%d&pageSize=%d", client.baseURL, page, pageSize)
	return doJSON[ListWorkoutsResponse](ctx, client.httpClient, http.MethodGet, url, nil)
}
