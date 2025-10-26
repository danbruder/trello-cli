# End-to-End Testing

This document describes the end-to-end (E2E) testing setup for the Trello CLI.

## Overview

The E2E test suite (`e2e_test.go`) exercises all major CLI commands against a live Trello API. This ensures that the CLI works correctly with real Trello boards, cards, lists, and other resources.

## Test Coverage

The E2E test covers the following command groups:

1. **Board Commands**
   - Create board
   - List boards
   - Get board details
   - Delete board (cleanup)

2. **List Commands**
   - Get lists from a board
   - Create list

3. **Card Commands**
   - Create card
   - List cards
   - Get card details
   - Move card between lists
   - Copy card
   - Archive card
   - Delete card

4. **Label Commands**
   - Create label
   - Get labels
   - Add label to card

5. **Checklist Commands**
   - Create checklist
   - Get checklists
   - Add checklist items

6. **Member Commands**
   - Get current member
   - Get board members

7. **Config Commands**
   - Load configuration

8. **Batch Commands**
   - Batch create cards

9. **Output Formats**
   - JSON format validation

## Running E2E Tests Locally

### Prerequisites

1. Trello API credentials:
   - Get your API key from: https://trello.com/app-key
   - Generate a token from the same page

2. Set environment variables:
   ```bash
   export TRELLO_API_KEY="your-api-key"
   export TRELLO_TOKEN="your-token"
   ```

### Running the Tests

Run the E2E test with:

```bash
go test -v -run TestE2EAllCommands ./...
```

Or with timeout and verbose output:

```bash
go test -v -timeout 30m -run TestE2EAllCommands ./...
```

### Test Behavior

- The test creates a temporary board with a unique timestamp-based name (e.g., `Test Board e2e-test-1234567890`)
- All test resources are created within this board
- The board is automatically deleted after the test completes (success or failure)
- The test skips if `TRELLO_API_KEY` or `TRELLO_TOKEN` are not set

## GitHub Actions

### Regular Tests

The `test.yml` workflow runs on every push and pull request:
- Tests against multiple Go versions (1.21, 1.22, 1.23)
- Runs unit tests with race detection
- Generates code coverage reports
- Uploads coverage to Codecov

### E2E Tests

The `e2e.yml` workflow runs E2E tests against live Trello API:

**Trigger conditions:**
- Manual trigger via GitHub Actions UI (workflow_dispatch)
- Weekly schedule (Mondays at 8 AM UTC)
- On push to main/master branches (when Go files change)

**Requirements:**
- GitHub repository secrets must be configured:
  - `TRELLO_API_KEY`: Your Trello API key
  - `TRELLO_TOKEN`: Your Trello token

### Setting Up GitHub Secrets

1. Go to your repository on GitHub
2. Navigate to Settings > Secrets and variables > Actions
3. Click "New repository secret"
4. Add the following secrets:
   - Name: `TRELLO_API_KEY`, Value: Your Trello API key
   - Name: `TRELLO_TOKEN`, Value: Your Trello token

## Test Output

The E2E test provides detailed logging:

```
=== RUN   TestE2EAllCommands
    e2e_test.go:33: Starting E2E test with ID: e2e-test-1698765432
=== RUN   TestE2EAllCommands/Board_Commands
=== RUN   TestE2EAllCommands/Board_Commands/Create_Board
    e2e_test.go:50: Created board: Test Board e2e-test-1698765432 (ID: 6543210abcdef)
=== RUN   TestE2EAllCommands/Board_Commands/List_Boards
    e2e_test.go:77: Found 5 boards
...
```

## Cleanup

- **Automatic**: The test automatically deletes the test board in a deferred cleanup function
- **Manual**: If a test fails and cleanup doesn't run, you can manually delete test boards from Trello
  - Look for boards named `Test Board e2e-test-*`

## Troubleshooting

### Test skipped with "requires Trello API credentials"

Set the required environment variables:
```bash
export TRELLO_API_KEY="your-api-key"
export TRELLO_TOKEN="your-token"
```

### API rate limiting

Trello has rate limits. If you encounter rate limit errors:
- Wait a few minutes between test runs
- The test is designed to be efficient and should stay within limits

### Authentication errors

Verify your credentials:
```bash
# Test with a simple API call
curl "https://api.trello.com/1/members/me?key=${TRELLO_API_KEY}&token=${TRELLO_TOKEN}"
```

### Leftover test boards

If cleanup fails, manually delete test boards:
1. Go to Trello
2. Find boards named `Test Board e2e-test-*`
3. Delete them manually

## Best Practices

1. **Don't commit credentials**: Never commit API keys or tokens to the repository
2. **Use secrets**: Always use GitHub Secrets for CI/CD credentials
3. **Monitor costs**: While Trello API is free, be aware of rate limits
4. **Review logs**: Check E2E test logs regularly to catch issues early

## Adding New Tests

To add new E2E test cases:

1. Add a new test run in `e2e_test.go`:
   ```go
   t.Run("Your New Test", func(t *testing.T) {
       // Your test code
   })
   ```

2. Follow the existing patterns:
   - Use `t.Fatalf()` for fatal errors
   - Use `t.Errorf()` for non-fatal assertions
   - Log important information with `t.Logf()`
   - Clean up resources if needed

3. Test locally before pushing:
   ```bash
   go test -v -run TestE2EAllCommands ./...
   ```

## Related Documentation

- [TESTING.md](TESTING.md) - General testing documentation
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Trello API Documentation](https://developer.atlassian.com/cloud/trello/rest/)
