package course

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository stores attempts only; it does not calculate progress or deduplicate IDs.
// Each operation is one SQL statement. Multiple calls are not a transaction.
type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) Add(ctx context.Context, attempt Attempt) error {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO attempts (attempt_id, learner_id, course_id, lesson_id, score)
 VALUES ($1, $2, $3, $4, $5)`,
		attempt.AttemptID,
		attempt.LearnerID,
		attempt.CourseID,
		attempt.LessonID,
		attempt.Score,
	)
	if err != nil {
		return fmt.Errorf("insert attempt: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, learnerID, courseID string) ([]Attempt, error) {
	rows, err := r.pool.Query(ctx, `SELECT attempt_id, learner_id, course_id, lesson_id, score
 FROM attempts WHERE learner_id = $1 AND course_id = $2 ORDER BY id`, learnerID, courseID)
	if err != nil {
		return nil, fmt.Errorf("list attempts: %w", err)
	}
	defer rows.Close()

	attempts := []Attempt{}

	for rows.Next() {
		var attempt Attempt
		if err := rows.Scan(
			&attempt.AttemptID,
			&attempt.LearnerID,
			&attempt.CourseID,
			&attempt.LessonID,
			&attempt.Score,
		); err != nil {
			return nil, fmt.Errorf("scan attempt: %w", err)
		}

		attempts = append(attempts, attempt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read attempts: %w", err)
	}

	return attempts, nil
}
