# Endless Quiz Flow - Database Design

## Overview

IQ Theory quizzes are **time-based with endless questions** - users answer as many note identification questions as possible within the time limit (30s, 60s, or 120s). The database is designed to support this endless question flow.

## Quiz Flow Design

### **1. Quiz Start**

```sql
-- User starts a quiz - only parameters are known, not total questions
INSERT INTO quiz_sessions (
    user_id,
    clef,
    duration_seconds,
    max_ledger_lines,
    status
) VALUES (
    'user-uuid',
    'treble',
    60,              -- 60 second time limit
    1,               -- up to 1 ledger line
    'in_progress'
);
-- total_questions starts at 0, will be incremented
-- time_taken_seconds is NULL until quiz completes
```

### **2. Endless Question Loop**

The frontend generates questions dynamically and the backend tracks each answer:

```sql
-- Question 1 answered
INSERT INTO quiz_answers (
    quiz_session_id,
    question_number,
    note_image,
    correct_note,
    user_answer,
    is_correct,
    time_taken_ms
) VALUES (
    'session-uuid',
    1,
    'C4.png',
    'C4',
    'C4',
    true,
    1500  -- 1.5 seconds to answer
);

-- Update session totals
UPDATE quiz_sessions SET
    total_questions = total_questions + 1,
    correct_answers = correct_answers + CASE WHEN true THEN 1 ELSE 0 END,
    score = score + CASE WHEN true THEN 1 ELSE 0 END  -- 1 point per correct answer
WHERE id = 'session-uuid';
```

This pattern repeats for each question until time runs out.

### **3. Quiz Completion (Time Up)**

```sql
-- When time limit reached, finalize the quiz
UPDATE quiz_sessions SET
    status = 'completed',
    completed_at = NOW(),
    time_taken_seconds = duration_seconds  -- Used full time allowance
WHERE id = 'session-uuid' AND status = 'in_progress';
```

## Database Schema Support

### **Key Schema Features for Endless Quizzes**:

1. **Dynamic Question Count**:

   ```sql
   total_questions INTEGER NOT NULL DEFAULT 0  -- Starts at 0, incremented per question
   ```

2. **Flexible Time Tracking**:

   ```sql
   time_taken_seconds INTEGER,  -- NULL during quiz, set on completion
   started_at TIMESTAMP,        -- Quiz start time
   completed_at TIMESTAMP       -- Quiz end time (when time limit reached)
   ```

3. **Real-time Accuracy**:

   ```sql
   accuracy_percentage DECIMAL(5,2) GENERATED ALWAYS AS (
       CASE WHEN total_questions > 0
            THEN (correct_answers::DECIMAL / total_questions) * 100
            ELSE 0 END
   ) STORED
   ```

4. **Question Sequence Tracking**:
   ```sql
   question_number INTEGER NOT NULL,  -- 1, 2, 3, 4... endless
   time_taken_ms INTEGER NOT NULL     -- Per-question timing
   ```

## Example Quiz Session Data

### **During Quiz (In Progress)**:

```json
{
  "id": "quiz-session-uuid",
  "user_id": "user-uuid",
  "clef": "treble",
  "duration_seconds": 60,
  "max_ledger_lines": 1,
  "total_questions": 15, // Answered 15 questions so far
  "correct_answers": 12, // Got 12 right
  "score": 12,
  "time_taken_seconds": null, // Still in progress
  "started_at": "2025-08-17T10:00:00Z",
  "completed_at": null, // Not finished yet
  "status": "in_progress",
  "accuracy_percentage": 80.0 // 12/15 = 80%
}
```

### **After Time Limit (Completed)**:

```json
{
  "id": "quiz-session-uuid",
  "user_id": "user-uuid",
  "clef": "treble",
  "duration_seconds": 60,
  "max_ledger_lines": 1,
  "total_questions": 23, // Final count: 23 questions in 60s
  "correct_answers": 18, // Final correct count
  "score": 18,
  "time_taken_seconds": 60, // Used full 60 seconds
  "started_at": "2025-08-17T10:00:00Z",
  "completed_at": "2025-08-17T10:01:00Z",
  "status": "completed",
  "accuracy_percentage": 78.26 // 18/23 = 78.26%
}
```

## API Endpoints for Endless Quizzes

### **Start Quiz**

```http
POST /api/quiz/start
{
  "clef": "treble",
  "duration_seconds": 60,
  "max_ledger_lines": 1
}

Response:
{
  "session_id": "uuid",
  "clef": "treble",
  "duration_seconds": 60,
  "started_at": "2025-08-17T10:00:00Z"
}
```

### **Submit Answer** (Called repeatedly)

```http
POST /api/quiz/{session_id}/answer
{
  "question_number": 1,
  "note_image": "C4.png",
  "correct_note": "C4",
  "user_answer": "C4",
  "time_taken_ms": 1500
}

Response:
{
  "is_correct": true,
  "current_score": 1,
  "total_questions": 1,
  "accuracy": 100.0,
  "time_remaining_seconds": 58
}
```

### **Get Next Question**

```http
GET /api/quiz/{session_id}/next-question

Response:
{
  "question_number": 2,
  "note_image": "F3.png",
  "possible_answers": ["F", "G", "E", "D"],  // Multiple choice options
  "time_remaining_seconds": 58
}
```

### **Quiz Auto-Complete** (When time runs out)

```http
PUT /api/quiz/{session_id}/complete
{
  "reason": "time_expired"
}

Response:
{
  "final_score": 18,
  "total_questions": 23,
  "accuracy_percentage": 78.26,
  "time_taken_seconds": 60,
  "quiz_completed": true
}
```

## Performance Considerations

### **Database Optimizations**:

1. **Batch Updates**: Update session totals every few questions rather than per question
2. **Indexing**: Index on (user_id, status) for active quiz lookup
3. **Connection Pooling**: Handle rapid question submission efficiently

### **Question Generation**:

- **Frontend Logic**: Generate questions client-side based on clef and ledger line limits
- **Image Caching**: Pre-load note images for smooth experience
- **Validation**: Server validates answers against known correct notes

### **Time Management**:

- **Client Timer**: JavaScript countdown timer on frontend
- **Server Validation**: Backend checks elapsed time before accepting answers
- **Auto-Complete**: Automatic quiz completion when time limit exceeded

## Summary

✅ **The updated database schema now fully supports endless time-based quizzes:**

- Dynamic question counting (starts at 0, increments with each answer)
- Flexible time tracking (NULL until completion)
- Real-time accuracy calculation
- Unlimited question storage per session
- Proper status management for in-progress vs completed quizzes

The design efficiently handles the endless question flow while maintaining detailed performance tracking for analytics and leaderboards.
