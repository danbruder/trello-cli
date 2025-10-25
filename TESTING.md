# Test Configuration

This file contains test configuration for the Trello CLI.

## Running Tests

### Unit Tests (No Trello API required)
```bash
go test ./internal/client/... -v
go test ./internal/formatter/... -v
go test ./internal/context/... -v
```

### Integration Tests (Requires Trello API credentials)
```bash
# Set environment variables
export TRELLO_API_KEY="your-api-key"
export TRELLO_TOKEN="your-token"

# Run integration tests
go test -v
```

### All Tests
```bash
go test ./... -v
```

## Test Coverage

To generate test coverage reports:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Structure

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions with real Trello API
- **Mock Tests**: Test with mock data when API is not available

## Test Data

Integration tests use real Trello data but are designed to be safe:
- Only read operations are performed
- No data is modified
- Tests skip gracefully if credentials are not provided
