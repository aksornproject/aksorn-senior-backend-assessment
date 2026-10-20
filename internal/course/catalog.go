package course

import (
	"fmt"
	"strings"
)

const CourseID = "course-001"

// FixedCourse returns a fresh configuration so callers cannot mutate shared state.
func FixedCourse() Course {
	return Course{
		ID: CourseID,
		Lessons: []Lesson{
			{ID: "introduction", Name: "Introduction", Required: true, PassingScore: 60},
			{ID: "practice", Name: "Practice", Required: true, PassingScore: 70},
			{ID: "challenge", Name: "Challenge", Required: false, PassingScore: 80},
		},
	}
}

// Validation is supplied. Candidates can assume valid inputs after this check.
func ValidateAttempt(a Attempt) error {
	ids := []string{a.AttemptID, a.LearnerID, a.CourseID, a.LessonID}
	for _, id := range ids {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("%w: all ids must be nonblank", ErrInvalidInput)
		}
	}

	if a.Score < 0 || a.Score > 100 {
		return fmt.Errorf("%w: score must be between 0 and 100", ErrInvalidInput)
	}

	if a.CourseID != CourseID {
		return ErrUnknownCourse
	}

	for _, lesson := range FixedCourse().Lessons {
		if lesson.ID == a.LessonID {
			return nil
		}
	}

	return ErrUnknownLesson
}
