# Efficient Quiz API Flow - Batch Processing Design

## Problem: API Request Per Question

### **Why Individual API Calls Are Bad:**

- **Network Overhead**: 50-200ms latency per question
- **Server Load**: 20+ requests per 60-second quiz
- **Database Pressure**: Constant INSERT operations
- **Poor UX**: Delays between questions break flow
- **Mobile Issues**: Higher latency, battery drain

### **Example Performance Impact:**

```
Traditional Approach (BAD):
- 25 questions in 60 seconds
- 25 API calls + start + end = 27 requests
- 100ms average latency = 2.7 seconds of waiting
- Database: 25 individual INSERT operations
```

## Recommended Solution: Batch Processing

### **Core Principles:**

1. **Frontend Question Generation**: Client generates questions locally
2. **Batch Answer Submission**: Submit 5-10 answers at once
3. **Periodic Sync**: Every 15-20 seconds or every 10 questions
4. **Final Flush**: Submit remaining data on completion

## Frontend Question Pool Structure

### **Bundled Note Data (Already in your app)**

Your existing quiz bank files are perfect for client-side generation:

```json
// /src/quizbank/treble.json
[
  {
    "note": "A",
    "imgs": ["A3", "A4", "A5"] // Different octaves/ledger positions
  },
  {
    "note": "B",
    "imgs": ["B3", "B4", "B5"]
  },
  {
    "note": "C",
    "imgs": ["C4", "C5", "C6"] // C6 would have ledger lines
  }
  // ... etc
]
```

### **Ledger Line Mapping**

```typescript
// Map note positions to ledger line counts
const LEDGER_LINE_MAP = {
  treble: {
    A3: 1,
    B3: 1, // Below staff
    C4: 0,
    D4: 0,
    E4: 0, // On staff
    F4: 0,
    G4: 0,
    A4: 0, // On staff
    B4: 0,
    C5: 0,
    D5: 0, // On staff
    E5: 0,
    F5: 0, // On staff
    G5: 1,
    A5: 1,
    B5: 2, // Above staff (ledger lines)
    C6: 2,
    D6: 3, // More ledger lines
  },
  bass: {
    // Similar mapping for bass clef positions
  },
  // ... etc for alto and tenor
};
```

## API Flow Design

### **1. Quiz Start (Single Request)**

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
  "started_at": "2025-08-17T10:00:00Z",
  "duration_seconds": 60,
  "clef": "treble",
  "max_ledger_lines": 1
}
```

**Note:** Questions are generated client-side from bundled note data, not retrieved from server.

### **2. Batch Answer Submission (Every 5-10 Questions)**

```http
POST /api/quiz/{session_id}/submit-batch
{
  "answers": [
    {
      "question_number": 1,
      "correct_note": "C4",
      "user_answer": "C4",
      "time_taken_ms": 1200,
      "answered_at": "2025-08-17T10:00:01.2Z"
    },
    {
      "question_number": 2,
      "correct_note": "G3",
      "user_answer": "A3",
      "time_taken_ms": 2100,
      "answered_at": "2025-08-17T10:00:03.3Z"
    }
    // ... 5-10 answers per batch
  ]
}

Response:
{
  "batch_processed": true,
  "current_score": 8,
  "total_questions": 10,
  "accuracy": 80.0,
  "time_remaining": 45
}
```

### **3. Quiz Completion (Final Submission)**

```http
POST /api/quiz/{session_id}/complete
{
  "final_answers": [
    // Any remaining unsubmitted answers
  ],
  "actual_time_used": 60,
  "completion_reason": "time_expired"
}

Response:
{
  "quiz_completed": true,
  "final_score": 18,
  "total_questions": 23,
  "accuracy_percentage": 78.26,
  "time_taken_seconds": 60,
  "rank_info": {
    "global_rank": 245,
    "group_rank": 12
  }
}
```

## Frontend Implementation Strategy

### **JavaScript Quiz Controller:**

```typescript
class QuizController {
  private pendingAnswers: QuizAnswer[] = [];
  private batchSize = 8; // Submit every 8 answers
  private syncInterval = 15000; // Or every 15 seconds
  private questionPool: NoteQuestion[];

  constructor() {
    // Load note data from bundled JSON files
    this.loadNoteData();
  }

  private async loadNoteData() {
    // These JSON files are bundled with the app
    const trebleNotes = await import("../assets/quizbank/treble.json");
    const bassNotes = await import("../assets/quizbank/bass.json");
    const altoNotes = await import("../assets/quizbank/alto.json");
    const tenorNotes = await import("../assets/quizbank/tenor.json");

    this.noteData = {
      treble: trebleNotes.default,
      bass: bassNotes.default,
      alto: altoNotes.default,
      tenor: tenorNotes.default,
    };
  }

  async startQuiz(params: StartQuizParams) {
    const response = await api.post("/quiz/start", params);
    this.sessionId = response.session_id;
    this.clef = params.clef;
    this.maxLedgerLines = params.max_ledger_lines;

    // Set up periodic sync
    this.syncTimer = setInterval(() => this.submitBatch(), this.syncInterval);
    this.generateNextQuestion();
  }

  submitAnswer(answer: QuizAnswer) {
    this.pendingAnswers.push(answer);
    this.generateNextQuestion(); // Immediate next question

    // Submit batch if we have enough answers
    if (this.pendingAnswers.length >= this.batchSize) {
      this.submitBatch();
    }
  }

  async submitBatch() {
    if (this.pendingAnswers.length === 0) return;

    const batch = this.pendingAnswers.splice(0); // Take all pending
    try {
      await api.post(`/quiz/${this.sessionId}/submit-batch`, {
        answers: batch,
      });
    } catch (error) {
      // Re-queue on failure for retry
      this.pendingAnswers.unshift(...batch);
    }
  }

  async completeQuiz(reason: string) {
    clearInterval(this.syncTimer);

    // Submit any remaining answers
    await api.post(`/quiz/${this.sessionId}/complete`, {
      final_answers: this.pendingAnswers,
      actual_time_used: this.getElapsedTime(),
      completion_reason: reason,
    });
  }

  generateNextQuestion() {
    // Client-side question generation - NO API CALL
    const availableNotes = this.getAvailableNotes();
    const question = this.selectRandomQuestion(availableNotes);
    this.displayQuestion(question);
  }

  private getAvailableNotes(): Note[] {
    const clefNotes = this.noteData[this.clef];

    // Filter notes based on ledger line limit
    return clefNotes.filter((note) => {
      return note.imgs.some((img) => {
        const ledgerCount = this.getLedgerLineCount(img);
        return ledgerCount <= this.maxLedgerLines;
      });
    });
  }

  private selectRandomQuestion(availableNotes: Note[]): QuizQuestion {
    const randomNote =
      availableNotes[Math.floor(Math.random() * availableNotes.length)];
    const validImages = randomNote.imgs.filter(
      (img) => this.getLedgerLineCount(img) <= this.maxLedgerLines
    );
    const randomImage =
      validImages[Math.floor(Math.random() * validImages.length)];

    return {
      noteImage: `${randomImage}.png`,
      correctNote: randomNote.note,
      correctAnswer: randomNote.note,
      options: this.generateMultipleChoiceOptions(randomNote.note),
    };
  }

  private getLedgerLineCount(noteImg: string): number {
    // Logic to determine ledger line count based on note position
    // e.g., "C6" in treble clef has 2 ledger lines above staff
    // This would be implemented based on your note naming convention
    const noteMap = {
      // Example for treble clef
      C6: 2,
      D6: 2,
      E6: 1,
      F6: 1,
      G6: 0,
      // ... etc for each clef
    };
    return noteMap[noteImg] || 0;
  }

  private generateMultipleChoiceOptions(correctNote: string): string[] {
    const allNotes = ["A", "B", "C", "D", "E", "F", "G"];
    const wrongOptions = allNotes.filter((note) => note !== correctNote);
    const selectedWrong = wrongOptions
      .sort(() => 0.5 - Math.random())
      .slice(0, 3);

    return [correctNote, ...selectedWrong].sort(() => 0.5 - Math.random());
  }
}
```

## Database Batch Processing

### **Efficient Batch INSERT:**

```sql
-- Single transaction for batch of answers
BEGIN;

-- Insert all answers in one statement
INSERT INTO quiz_answers (
    quiz_session_id, question_number,
    correct_note, user_answer, is_correct,
    time_taken_ms, answered_at
) VALUES
    ('session-uuid', 1, 'C4', 'C4', true, 1200, '2025-08-17T10:00:01.2Z'),
    ('session-uuid', 2, 'G3', 'A3', false, 2100, '2025-08-17T10:00:03.3Z'),
    -- ... up to 10 answers per batch
;

-- Update session totals in single operation
UPDATE quiz_sessions SET
    total_questions = total_questions + 8,    -- Batch size
    correct_answers = correct_answers + 6,    -- Count of correct in batch
    score = score + 6
WHERE id = 'session-uuid';

COMMIT;
```

## Performance Comparison

### **Traditional (Per Question API):**

```
Quiz Start: 1 request (with question pool download)
Network Requests: 25 questions = 25 requests
Total Latency: 26 × 100ms = 2.6 seconds
Database Operations: 25 individual INSERTs
User Experience: Lag between every question
```

### **Frontend Question Pool + Batch Processing:**

```
Quiz Start: 1 request (session only, no question data)
Question Generation: Client-side (instant, no network)
Network Requests: 3-4 batches = 4 requests total
Total Latency: 4 × 100ms = 0.4 seconds
Database Operations: 3-4 batch INSERTs
User Experience: Instant questions, seamless flow
```

**Result: 85% reduction in network overhead + instant question generation!**

## Fallback Strategies

### **Network Issues:**

```typescript
// Retry failed batches
private async submitBatchWithRetry(batch: QuizAnswer[], maxRetries = 3) {
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      await this.submitBatch(batch);
      return;
    } catch (error) {
      if (attempt === maxRetries - 1) {
        // Store locally for later submission
        this.storeOffline(batch);
      }
      await this.delay(1000 * Math.pow(2, attempt)); // Exponential backoff
    }
  }
}
```

### **Offline Support:**

```typescript
// Store answers locally if network fails
private storeOffline(answers: QuizAnswer[]) {
  const offline = JSON.parse(localStorage.getItem('offline_answers') || '[]');
  offline.push(...answers);
  localStorage.setItem('offline_answers', JSON.stringify(offline));
}

// Submit when network returns
private async submitOfflineAnswers() {
  const offline = JSON.parse(localStorage.getItem('offline_answers') || '[]');
  if (offline.length > 0) {
    await this.submitBatch(offline);
    localStorage.removeItem('offline_answers');
  }
}
```

## Summary

**❌ Don't:** Send API request per question

- High latency, poor UX, server overload

**❌ Don't:** Download question pool from backend

- Larger API responses, network delay at quiz start

**✅ Do:** Frontend question pool + batch processing

- **Bundle note data** with the app (you already have this!)
- **Generate questions client-side** from existing JSON files
- **Submit answers in batches** of 5-10 questions
- **Periodic sync** every 15-20 seconds
- **Final submission** on completion

**Benefits:**

- ✅ **Instant quiz start** - no network delay
- ✅ **Seamless question flow** - no API calls during quiz
- ✅ **85% fewer API requests** - better performance
- ✅ **Offline capable** - works without network
- ✅ **Smaller app size** - note data is tiny JSON files

**Perfect for IQ Theory because:**

- Note identification is deterministic (C4 is always C4)
- Limited, well-defined question set per clef
- Visual assets already bundled with app
- User experience prioritizes speed and responsiveness

This approach provides the best possible quiz experience with minimal server load.
