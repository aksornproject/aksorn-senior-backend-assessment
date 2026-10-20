package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/aksornproject/aksorn-senior-backend-assessment/internal/course"
	"github.com/aksornproject/aksorn-senior-backend-assessment/internal/infra"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	connections, err := infra.Open(context.Background())
	if err != nil {
		return err
	}
	defer connections.Close()

	fmt.Println("PASS PostgreSQL, supplied schema, and Redis connections")

	setupID := fmt.Sprintf("setup-%d", time.Now().UnixNano())
	handler := course.NewHandler(course.NewService(course.NewRepository(connections.Postgres)))

	checks := []struct {
		name     string
		method   string
		path     string
		body     string
		statuses []int
	}{
		{name: "health", method: http.MethodGet, path: "/healthz", statuses: []int{200}},
		{
			name:     "JSON validation",
			method:   http.MethodPost,
			path:     "/attempts",
			body:     "{",
			statuses: []int{400},
		},
		{
			name:     "score validation",
			method:   http.MethodPost,
			path:     "/attempts",
			body:     `{"attempt_id":"setup","learner_id":"setup","course_id":"course-001","lesson_id":"introduction","score":101}`,
			statuses: []int{400},
		},
		{
			name:     "attempt endpoint",
			method:   http.MethodPost,
			path:     "/attempts",
			body:     `{"attempt_id":"setup","learner_id":"setup","course_id":"course-001","lesson_id":"introduction","score":60}`,
			statuses: []int{501, 201},
		},
		{
			name: "progress endpoint", method: http.MethodGet,
			path: "/learners/setup/courses/course-001/progress", statuses: []int{501, 200},
		},
	}
	for _, check := range checks {
		request := httptest.NewRequest(
			check.method,
			strings.ReplaceAll(check.path, "setup", setupID),
			strings.NewReader(strings.ReplaceAll(check.body, "setup", setupID)),
		)
		request.Header.Set("Content-Type", "application/json")

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		ok := false

		for _, status := range check.statuses {
			if status == response.Code {
				ok = true
			}
		}

		if !ok || !json.Valid(response.Body.Bytes()) {
			return fmt.Errorf(
				"%s failed: HTTP %d %s",
				check.name,
				response.Code,
				response.Body.String(),
			)
		}

		fmt.Printf("PASS %s (HTTP %d)\n", check.name, response.Code)
	}

	fmt.Println("Setup verified. HTTP 501 is expected until the candidate TODOs are implemented.")
	fmt.Println("This checks starter wiring, not business-rule correctness.")

	return nil
}
