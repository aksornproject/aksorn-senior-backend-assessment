package course

import "errors"

var (
	ErrNotImplemented = errors.New("course: candidate implementation required")
	ErrInvalidInput   = errors.New("course: invalid input")
	ErrUnknownCourse  = errors.New("course: unknown course")
	ErrUnknownLesson  = errors.New("course: unknown lesson")
	ErrConflict       = errors.New("course: conflicting attempt id")
)

type Attempt struct {
	AttemptID string `json:"attempt_id"`
	LearnerID string `json:"learner_id"`
	CourseID  string `json:"course_id"`
	LessonID  string `json:"lesson_id"`
	Score     int    `json:"score"`
}

type Lesson struct {
	ID           string
	Name         string
	Required     bool
	PassingScore int
}

type Course struct {
	ID      string
	Lessons []Lesson
}

type LessonProgress struct {
	LessonID  string `json:"lesson_id"`
	BestScore *int   `json:"best_score"`
	Passed    bool   `json:"passed"`
}

type Progress struct {
	LearnerID       string           `json:"learner_id"`
	CourseID        string           `json:"course_id"`
	Lessons         []LessonProgress `json:"lessons"`
	ProgressPercent int              `json:"progress_percent"`
	Completed       bool             `json:"completed"`
}
