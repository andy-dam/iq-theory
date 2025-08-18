# Schema Optimization: Removing note_image Field

## Question Analysis

**Should `note_image` be stored in the `quiz_answers` table?**

**Answer: No** - The `note_image` field is redundant and should be removed.

## Why `note_image` is Unnecessary

### **1. Redundant Data**

```sql
-- BEFORE (redundant)
note_image VARCHAR(100) NOT NULL, -- "A4.png"
correct_note VARCHAR(10) NOT NULL, -- "A4"

-- AFTER (sufficient)
correct_note VARCHAR(10) NOT NULL, -- "A4" tells us everything
```

### **2. Frontend Controls Images**

With client-side question generation:

- Frontend selects which image to show (`A4.png`, `A4_alt.png`, etc.)
- Server doesn't need to know which specific image was displayed
- `correct_note = "A4"` is sufficient to identify the question

### **3. Deterministic Mapping**

```typescript
// Frontend can always reconstruct what was shown
const noteToImage = (note: string, clef: string) => {
  return `${clef}/${note}.png`; // A4 → treble/A4.png
};
```

### **4. Storage Efficiency**

```
Per quiz session with 25 questions:
- BEFORE: 25 × 100 bytes = 2.5 KB extra storage
- AFTER: 25 × 0 bytes = 0 KB
- SAVINGS: 100% reduction in image field storage

With 1M quiz sessions: 2.5 GB savings
```

## Updated Schema

### **Database Table:**

```sql
CREATE TABLE quiz_answers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_session_id UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    question_number INTEGER NOT NULL,
    correct_note VARCHAR(10) NOT NULL, -- Sufficient identifier
    user_answer VARCHAR(10),
    is_correct BOOLEAN NOT NULL,
    time_taken_ms INTEGER NOT NULL,
    answered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(quiz_session_id, question_number)
);
```

### **API Request:**

```json
// Batch answer submission (no note_image needed)
{
  "answers": [
    {
      "question_number": 1,
      "correct_note": "A4", // Frontend knows this was A4
      "user_answer": "A4", // User's guess
      "time_taken_ms": 1200,
      "answered_at": "2025-08-17T10:00:01.2Z"
    }
  ]
}
```

### **Go Models:**

```go
type QuizAnswer struct {
    ID             uuid.UUID `json:"id" db:"id"`
    QuizSessionID  uuid.UUID `json:"quiz_session_id" db:"quiz_session_id"`
    QuestionNumber int       `json:"question_number" db:"question_number"`
    CorrectNote    string    `json:"correct_note" db:"correct_note"`  // Sufficient
    UserAnswer     *string   `json:"user_answer" db:"user_answer"`
    IsCorrect      bool      `json:"is_correct" db:"is_correct"`
    TimeTakenMs    int       `json:"time_taken_ms" db:"time_taken_ms"`
    AnsweredAt     time.Time `json:"answered_at" db:"answered_at"`
}
```

## What We Don't Lose

### **Analytics Still Possible:**

```sql
-- Most difficult notes
SELECT correct_note,
       COUNT(*) as total_attempts,
       AVG(CASE WHEN is_correct THEN 1.0 ELSE 0.0 END) as accuracy
FROM quiz_answers qa
JOIN quiz_sessions qs ON qa.quiz_session_id = qs.id
WHERE qs.clef = 'treble'
GROUP BY correct_note
ORDER BY accuracy ASC;

-- Performance by clef and ledger lines
SELECT qs.clef, qs.max_ledger_lines,
       AVG(qa.time_taken_ms) as avg_time,
       AVG(CASE WHEN qa.is_correct THEN 1.0 ELSE 0.0 END) as accuracy
FROM quiz_answers qa
JOIN quiz_sessions qs ON qa.quiz_session_id = qs.id
GROUP BY qs.clef, qs.max_ledger_lines;
```

### **Debugging Still Possible:**

```sql
-- Find problematic answers
SELECT * FROM quiz_answers
WHERE correct_note = 'A4' AND is_correct = false
ORDER BY answered_at DESC;
```

## Migration Strategy

### **For Existing Data:**

```sql
-- If you already have note_image data, you can safely drop it
ALTER TABLE quiz_answers DROP COLUMN note_image;
```

### **For New Development:**

- Remove `note_image` from table creation
- Update API endpoints to not expect `note_image`
- Frontend continues working as-is (doesn't send image info)

## Summary

✅ **Benefits of removing `note_image`:**

- Reduces storage by 100 bytes per answer
- Simplifies API contracts
- Eliminates redundant data
- Frontend already controls image selection
- No loss of analytical capability

❌ **No significant downsides:**

- `correct_note` provides all necessary information
- Analytics queries remain fully functional
- Debugging capabilities preserved

**Recommendation:** Remove `note_image` field entirely for a cleaner, more efficient schema.
