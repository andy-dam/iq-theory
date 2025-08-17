# Tests

This directory contains your test files.

## Types of Tests

- **Unit tests**: Test individual functions/methods (`*_test.go` files alongside your code)
- **Integration tests**: Test how components work together
- **End-to-end tests**: Test complete user workflows

## Testing Tools

- Built-in `testing` package
- `testify` for assertions and mocks
- `httptest` for HTTP handler testing
- Database testing with test containers or in-memory databases

## Test Structure

```
tests/
├── integration/        # Integration tests
├── e2e/               # End-to-end tests
└── fixtures/          # Test data files
```
