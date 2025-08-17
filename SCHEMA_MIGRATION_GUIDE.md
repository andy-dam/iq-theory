# Schema Migration Guide - Efficient Quiz Configuration

## Overview

This document outlines the migration from the original quiz configuration approach to the new efficient, normalized schema design.

## What Changed

### Before (Inefficient)

- **48 pre-stored quiz configuration records** with redundant data
- Each configuration stored as a separate row: `('Treble 30s - No Ledger', 'treble', 30, 0)`
- Hard to maintain and extend
- Storage waste due to string repetition

### After (Efficient)

- **11 normalized parameter records** (4 clefs + 3 durations + 4 ledger options)
- **Dynamic configuration generation** via views and API
- **78% reduction** in stored records
- Easy to add new options

## Database Schema Changes

### Removed Tables

```sql
-- REMOVED: quiz_configurations table
-- No longer needed - replaced with parameter tables
```

### Added Tables

```sql
-- NEW: Parameter tables for normalization
CREATE TABLE clef_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE duration_options (
    id SERIAL PRIMARY KEY,
    duration_seconds INTEGER UNIQUE NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE ledger_line_options (
    id SERIAL PRIMARY KEY,
    max_lines INTEGER UNIQUE NOT NULL,
    display_name VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true
);
```

### Modified Tables

```sql
-- MODIFIED: quiz_sessions table
-- OLD: quiz_configuration_id UUID REFERENCES quiz_configurations(id)
-- NEW: Individual parameter fields with foreign key constraints

ALTER TABLE quiz_sessions
DROP COLUMN quiz_configuration_id,
ADD COLUMN clef VARCHAR(20) NOT NULL,
ADD COLUMN duration_seconds INTEGER NOT NULL,
ADD COLUMN max_ledger_lines INTEGER NOT NULL,
ADD FOREIGN KEY (clef) REFERENCES clef_types(name),
ADD FOREIGN KEY (duration_seconds) REFERENCES duration_options(duration_seconds),
ADD FOREIGN KEY (max_ledger_lines) REFERENCES ledger_line_options(max_lines);
```

### Added Views

```sql
-- NEW: Dynamic configuration view
CREATE VIEW available_quiz_configurations AS
SELECT
    CONCAT(ct.display_name, ' - ', do.display_name, ' - ', llo.display_name) as configuration_name,
    ct.name as clef,
    do.duration_seconds,
    llo.max_lines as max_ledger_lines,
    ct.is_active AND do.is_active AND llo.is_active as is_available
FROM clef_types ct
CROSS JOIN duration_options do
CROSS JOIN ledger_line_options llo;
```

## API Changes

### Before (Configuration ID Based)

```go
// Start Quiz Request
POST /api/quiz/start
{
    "quiz_configuration_id": "uuid-here"
}

// Get Configurations
GET /api/quiz-configurations
// Returns: Array of 48 pre-stored configurations
```

### After (Parameter Based)

```go
// Start Quiz Request
POST /api/quiz/start
{
    "clef": "treble",
    "duration_seconds": 60,
    "max_ledger_lines": 2
}

// Get Available Options
GET /api/quiz-options
// Returns: {clefs: [...], durations: [...], ledger_lines: [...]}

// Get All Combinations (if needed)
GET /api/quiz-configurations
// Returns: Dynamic view of all 48 combinations
```

## Frontend Changes Required

### Quiz Configuration Selection

```typescript
// Before: Select from dropdown of 48 pre-defined configurations
interface QuizConfig {
  id: string;
  name: string;
  clef: string;
  duration_seconds: number;
  max_ledger_lines: number;
}

// After: Build configuration from separate parameter selections
interface QuizOptions {
  clefs: ClefType[];
  durations: DurationOption[];
  ledger_lines: LedgerLineOption[];
}

interface StartQuizParams {
  clef: string;
  duration_seconds: number;
  max_ledger_lines: number;
}
```

### Component Updates

```tsx
// Before: Single dropdown with 48 options
<select name="quiz_config">
    {configurations.map(config =>
        <option value={config.id}>{config.name}</option>
    )}
</select>

// After: Three separate dropdowns
<select name="clef">
    {options.clefs.map(clef =>
        <option value={clef.name}>{clef.display_name}</option>
    )}
</select>
<select name="duration">
    {options.durations.map(duration =>
        <option value={duration.duration_seconds}>{duration.display_name}</option>
    )}
</select>
<select name="ledger_lines">
    {options.ledger_lines.map(option =>
        <option value={option.max_lines}>{option.display_name}</option>
    )}
</select>
```

## Migration Benefits

### Storage Efficiency

- **Before**: 48 records × ~100 bytes = ~4.8KB
- **After**: 11 records × ~50 bytes = ~550 bytes
- **Savings**: 88% reduction in configuration storage

### Maintenance Efficiency

- **Before**: Add new duration → Insert 16 new configuration records
- **After**: Add new duration → Insert 1 duration option record
- **Benefit**: Automatic generation of all new combinations

### Query Performance

- **Before**: JOIN on large quiz_configurations table
- **After**: Direct filtering on indexed parameter columns
- **Benefit**: Faster leaderboard and filtering queries

### Flexibility

- **Before**: Fixed set of 48 combinations
- **After**: Easy to enable/disable parameter types, add new options
- **Benefit**: Runtime configuration management

## Testing Considerations

### Database Testing

```sql
-- Verify all combinations are available
SELECT COUNT(*) FROM available_quiz_configurations; -- Should return 48

-- Test parameter constraints
INSERT INTO quiz_sessions (clef, duration_seconds, max_ledger_lines)
VALUES ('invalid_clef', 60, 1); -- Should fail with foreign key constraint

-- Test leaderboard performance
EXPLAIN ANALYZE SELECT * FROM leaderboards
WHERE clef = 'treble' AND duration_seconds = 60;
```

### API Testing

```go
// Test parameter validation
func TestStartQuizValidation(t *testing.T) {
    // Invalid clef should fail
    req := StartQuizRequest{Clef: "invalid", DurationSeconds: 60, MaxLedgerLines: 1}
    assert.Error(t, validateStartQuizRequest(req))

    // Valid parameters should succeed
    req = StartQuizRequest{Clef: "treble", DurationSeconds: 60, MaxLedgerLines: 1}
    assert.NoError(t, validateStartQuizRequest(req))
}
```

### Frontend Testing

```typescript
// Test quiz options loading
test("loads quiz options correctly", async () => {
  const options = await fetchQuizOptions();
  expect(options.clefs).toHaveLength(4);
  expect(options.durations).toHaveLength(3);
  expect(options.ledger_lines).toHaveLength(4);
});

// Test quiz start with parameters
test("starts quiz with parameters", async () => {
  const params = { clef: "treble", duration_seconds: 60, max_ledger_lines: 2 };
  const session = await startQuiz(params);
  expect(session.clef).toBe("treble");
  expect(session.duration_seconds).toBe(60);
});
```

## Rollback Plan

If needed, the old approach can be restored by:

1. **Recreate quiz_configurations table**
2. **Populate with 48 combinations** from parameter tables
3. **Add quiz_configuration_id back to quiz_sessions**
4. **Revert API endpoints** to use configuration IDs
5. **Update frontend** to use single dropdown

However, the new approach is recommended for long-term maintainability and performance.
