# Technical Note

## 1. Design decision

Describe one decision in your implementation, why you chose it, and its tradeoff.

[Your answer]

## 2. Design questions

Brief answers are sufficient. No additional implementation is required.

### Retried requests

How would you enforce globally unique attempt IDs, return the original result for an identical replay, and return 409 for different content under the same ID? Consider simultaneous requests with the same ID.

[Your answer]

### Concurrent requests

Two attempts with distinct IDs and scores of 90 and 70 arrive simultaneously for the same learner and lesson.

Can your implementation lose the higher score? What could a simultaneous progress read observe? What should a new read return after both writes have committed? Is any protection needed?

[Your answer]

### Redis

Suppose attempt history grows and progress reads become slow.

Would you use Redis? What would you store, how would it stay consistent with PostgreSQL, and what would happen if Redis were unavailable or its data were lost? An alternative approach is welcome.

[Your answer]

## 3. Limitations and production readiness

List any limitations or unfinished work, and what you would verify before production use.

[Your answer]

## 4. Time and AI tools

Time spent: [Actual time spent, up to four hours]

AI tools used: [None, or tools used and what they helped with]
