package course

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type attemptRequest struct {
	AttemptID string `json:"attempt_id"`
	LearnerID string `json:"learner_id"`
	CourseID  string `json:"course_id"`
	LessonID  string `json:"lesson_id"`
	Score     *int   `json:"score"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(service *Service) http.Handler {
	mux := chi.NewRouter()
	mux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.Post("/attempts", func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			writeJSON(
				w,
				http.StatusUnsupportedMediaType,
				errorResponse{Error: "use application/json"},
			)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var input attemptRequest
		if err := decoder.Decode(&input); err != nil {
			writeDecodeError(w, err)
			return
		}

		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			writeDecodeError(w, err)
			return
		}

		if input.Score == nil {
			writeJSON(
				w,
				http.StatusBadRequest,
				errorResponse{Error: "score is required and must be an integer"},
			)
			return
		}

		attempt := Attempt{
			AttemptID: input.AttemptID,
			LearnerID: input.LearnerID,
			CourseID:  input.CourseID,
			LessonID:  input.LessonID,
			Score:     *input.Score,
		}
		if err := service.RecordAttempt(r.Context(), attempt); err != nil {
			writeServiceError(w, err)
			return
		}

		writeJSON(w, http.StatusCreated, attempt)
	})
	mux.Get(
		"/learners/{learner_id}/courses/{course_id}/progress",
		func(w http.ResponseWriter, r *http.Request) {
			progress, err := service.Progress(
				r.Context(),
				chi.URLParam(r, "learner_id"),
				chi.URLParam(r, "course_id"),
			)
			if err != nil {
				writeServiceError(w, err)
				return
			}

			writeJSON(w, http.StatusOK, progress)
		},
	)

	return mux
}

func writeDecodeError(w http.ResponseWriter, err error) {
	var tooLarge *http.MaxBytesError
	if errors.As(err, &tooLarge) {
		writeJSON(
			w,
			http.StatusRequestEntityTooLarge,
			errorResponse{Error: "request body is too large"},
		)
		return
	}

	writeJSON(
		w,
		http.StatusBadRequest,
		errorResponse{Error: "provide one valid JSON object with known fields"},
	)
}

func writeServiceError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "internal error"

	switch {
	case errors.Is(err, ErrNotImplemented):
		status, message = http.StatusNotImplemented, "candidate implementation required in internal/course/service.go"
	case errors.Is(err, ErrInvalidInput):
		status, message = http.StatusBadRequest, err.Error()
	case errors.Is(err, ErrUnknownCourse), errors.Is(err, ErrUnknownLesson):
		status, message = http.StatusNotFound, err.Error()
	case errors.Is(err, ErrConflict):
		status, message = http.StatusConflict, err.Error()
	default:
		slog.Error("application request failed", "error", err)
	}

	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("response encoding failed", "error", err)
	}
}
