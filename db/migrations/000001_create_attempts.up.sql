-- Supplied exercise storage. Attempt-ID deduplication is an optional extension.
CREATE TABLE attempts (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    attempt_id text NOT NULL,
    learner_id text NOT NULL,
    course_id text NOT NULL,
    lesson_id text NOT NULL,
    score integer NOT NULL CHECK (score BETWEEN 0 AND 100)
);
CREATE INDEX attempts_learner_course_idx ON attempts (learner_id, course_id);
