# Fictional learner progress assessment

## Purpose and time limit

This is an assessment-only exercise. All requirements are invented and unrelated to Aksorn's product roadmap. Your work will not be incorporated into Aksorn products.

Implement a small part of a fictional learning platform in Go. We assess how you understand business rules, preserve correctness, test behavior, and explain decisions. We do not assess how many features you can add.

Verify setup using the README before starting. The **four hours include reading this assignment, implementation, tests, and a short technical note**. Stop at four hours and identify unfinished work. Incomplete submissions are reviewed fairly, but core correctness still matters. The design questions require brief written answers only; no additional implementation is required.

## Supplied application

The server, routes, JSON parsing, validation, response encoding, fixed course, and PostgreSQL repository already work; Chi routing, pgxpool connections, Redis connectivity, Docker Compose, and golang-migrate schema migrations are supplied. The two application methods in `internal/course/service.go` deliberately return `ErrNotImplemented`. HTTP 501 is their expected initial response. Replace those returns with your implementation.

The repository supplies `Add(ctx, attempt)` and `List(ctx, learnerID, courseID)`. It preserves all inserted attempts and returns records filtered by learner and course. It does not calculate progress or deduplicate IDs. Each method executes one SQL statement. Separate calls are not a transaction.

Assume **sequential requests**, fixed configuration, and trusted internal callers. Authentication, prerequisites, idempotency, concurrent-request safety, a frontend, and additional infrastructure are outside the required implementation. Keep Chi and pgxpool. Plain SQL, SQLC, SQLX, and Bun are allowed; no other ORMs, including GORM or Ent. Use distinct attempt IDs in baseline examples and tests. Retry semantics are outside the required implementation and are discussed in the design questions.

Persistence and Redis connection setup are provided. Calculating progress directly from stored attempts is sufficient for this assignment. You do not need to implement caching, other Redis behavior, or a separate progress aggregate. Redis is included for the design discussion below; using it in the implementation earns no credit by itself.

## Course and input

There is one course, `course-001`:

| Lesson ID | Name | Required | Passing score |
| --- | --- | --- | --- |
| introduction | Introduction | Yes | 60 |
| practice | Practice | Yes | 70 |
| challenge | Challenge | No | 80 |

An attempt contains `attempt_id`, `learner_id`, `course_id`, `lesson_id`, and `score`. IDs are nonblank, case-sensitive strings. Learners need no registration. The score is a required JSON integer from 0 through 100. Zero is a valid score; a missing or null score is invalid. The starter already rejects malformed JSON, unknown fields, invalid scores, unknown courses, and unknown lessons. There are no prerequisites.

## Required implementation

### Record an attempt

Complete `Service.RecordAttempt` to store an accepted attempt with the supplied repository. Keep the provided validation. Support multiple attempts by the same learner for the same lesson. Once storage succeeds, the supplied HTTP handler returns 201 and echoes the accepted attempt.

### Calculate progress

Complete `Service.Progress`. Return all three lessons, including ones with no attempts, and apply these rules:

1. Each lesson's best score is the highest accepted score for that learner and course.
2. A score equal to or greater than the lesson's passing score passes.
3. A later lower score must not remove an earlier passing result or lower the best score.
4. Only required lessons count toward overall progress: neither passed is 0%, one passed is 50%, both passed is 100%.
5. The course is complete only when both required lessons have passed.
6. The optional lesson still has its own best score and pass status, but does not affect overall progress or completion.
7. One learner's attempts must not affect another learner's results.
8. With no attempts, `best_score` is `null` and `passed` is `false`. A real score of 0 must appear as 0, not null.

The lesson array may be in any order. Preserve the response fields below. For a learner with only one accepted score of 60 on Introduction:

```json
{
  "learner_id": "learner-001",
  "course_id": "course-001",
  "lessons": [
    {"lesson_id": "introduction", "best_score": 60, "passed": true},
    {"lesson_id": "practice", "best_score": null, "passed": false},
    {"lesson_id": "challenge", "best_score": null, "passed": false}
  ],
  "progress_percent": 50,
  "completed": false
}
```

## HTTP contract

`POST /attempts` accepts `application/json` and returns 201 with the attempt after successful recording. `GET /learners/{learner_id}/courses/{course_id}/progress` returns 200 with the progress object. A valid learner with no attempts receives 200 and empty progress, not 404. `GET /healthz` returns 200 with `{"status":"ok"}`.

The supplied error envelope is `{"error":"description"}`. Invalid input returns 400; unknown course or lesson returns 404; a body over 64 KiB returns 413; an unsupported content type returns 415. Missing application logic returns 501. Unknown routes and unsupported methods use Chi 404/405 responses. You do not need to change HTTP plumbing.

## Automated tests

Write focused tests for multiple attempts and best scores, passing thresholds, retained passing status, required versus optional lessons, course completion, and learner isolation. Include boundaries that would catch a plausible mistake. Application tests are sufficient; HTTP tests are optional. There is no coverage-percentage requirement or required internal architecture.

## Design questions

Answer all three questions briefly in `TECHNICAL_NOTE.md`. **Discussion only; implementation is not required.** Base your answers on your submitted design and state any assumptions.

1. **Retried requests:** Suppose attempt IDs must be globally unique. An identical replay should return the original accepted result without storing another attempt; different content under the same ID should return 409 without changing data. How would you enforce this, including if two requests with the same ID arrive simultaneously?
2. **Concurrent requests:** Two attempts with distinct IDs and scores of 90 and 70 arrive simultaneously for the same learner and lesson. Can your implementation lose the higher score? What could a simultaneous progress read observe, and what should a new read return after both writes have committed? Explain whether any protection is needed.
3. **Redis:** Suppose the attempt history grows and progress reads become slow. Would you use Redis? Explain what you would store, how it would stay consistent with PostgreSQL, and what would happen if Redis were unavailable or its data were lost. You may propose an alternative if Redis is not appropriate.

We assess reasoning and tradeoffs, not the number of technologies proposed. Identifying that a particular risk does not exist in your design is a valid answer. Do not implement these scenarios within the assessment.

## Technical note and submission

complete `TECHNICAL_NOTE.md`, containing:

1. One design decision and its tradeoff.
2. Brief answers to all three design questions above.
3. Limitations or unfinished work, and what you would verify before production use.
4. Time spent and any AI tools used.

AI tools are permitted; you must understand and be able to explain and change the submitted code.

Follow this submission workflow:

1. Fork the assessment repository to your personal GitHub account.
2. Develop the assignment on your fork.
3. Commit changes progressively as you work.
4. Push your work to GitHub.
5. Ensure that the hiring team can access your repository, either publicly or through granted private access.
6. Use **Reply All** on the recruitment email and include your repository link.

Stop after four hours. Include your code, tests, and `TECHNICAL_NOTE.md`. Exclude binaries and local caches. You may include additional notes in your reply about tradeoffs, limitations, unfinished work, or special setup instructions. These notes do not replace `TECHNICAL_NOTE.md`.

See the [README setup instructions](README.md#prerequisites-and-setup) before starting. If GitHub prevents you from forking or accessing the repository, contact the hiring team before starting the timer. No diagrams, extensive documentation, deployment, or services beyond the supplied PostgreSQL and Redis are requested.

## Evaluation

The baseline is a correct, proportionate implementation of the required rules, meaningful tests, readable Go, and honest explanation of limitations. The technical note and a brief interview follow-up assess risk reasoning and tradeoffs. You do not need to redesign the supplied persistence or implement concurrent-request safety to pass. A small omission can be discussed; broad correctness gaps are not offset by extra features. Documenting a defect does not make it correct. The design questions assess reasoning only; additional implementation does not earn credit by itself.
