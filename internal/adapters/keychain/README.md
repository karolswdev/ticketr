# Keychain Credential Store

## Overview

This package provides a production-ready implementation of the `ports.CredentialStore` interface using OS-level keychain integration via the `github.com/zalando/go-keyring` library.

## Features

- **Cross-platform Support**: Works with macOS Keychain, Windows Credential Manager, and Linux Secret Service API
- **Secure Storage**: Leverages OS-level encryption for credential protection
- **Thread-safe**: All operations are thread-safe via OS synchronization mechanisms
- **Comprehensive Error Handling**: Proper error messages with context
- **Security Best Practices**: Never logs credentials or sensitive data

## Architecture

This adapter follows the hexagonal architecture pattern:

```
ports.CredentialStore (interface)
    ↓
KeychainStore (adapter implementation)
    ↓
github.com/zalando/go-keyring (library)
    ↓
OS Keychain (macOS/Windows/Linux)
```

## Usage

```go
import "github.com/karolswdev/ticktr/internal/adapters/keychain"

// Create a new keychain store
store := keychain.NewKeychainStore()

// Store credentials
config := domain.WorkspaceConfig{
    JiraURL:    "https://example.atlassian.net",
    ProjectKey: "PROJ",
    Username:   "user@example.com",
    APIToken:   "secret-token",
}

ref, err := store.Store("workspace-id", config)
if err != nil {
    log.Fatal(err)
}

// Retrieve credentials
retrievedConfig, err := store.Retrieve(ref)
if err != nil {
    log.Fatal(err)
}

// List all credentials
refs, err := store.List()
if err != nil {
    log.Fatal(err)
}

// Delete credentials
err = store.Delete(ref)
if err != nil {
    log.Fatal(err)
}
```

## Implementation Details

### Service Name
All credentials are stored under the service name `"ticketr"` in the OS keychain.

### Credential Reference
Each stored credential returns a `domain.CredentialRef` containing:
- `KeychainID`: The workspace ID (used as the keychain account identifier)
- `ServiceID`: Always "ticketr"

### Storage Format
Credentials are stored as JSON-encoded `domain.WorkspaceConfig` structures in the OS keychain.

### Credential List Management
The adapter maintains an internal list of all stored credential IDs using a special keychain entry with ID `"ticketr-credential-list"`. This enables the `List()` operation without requiring OS-specific keychain enumeration APIs.

## Testing

### Unit Tests
The package includes comprehensive unit tests that validate:
- Input validation logic
- Error detection (cross-platform "not found" error patterns)
- Domain validation
- Interface compliance

These tests run without requiring actual keychain access.

### Integration Tests
Full integration tests are included that test against the real OS keychain:
- Store/Retrieve round-trip
- Delete operations
- List operations
- Concurrent access
- Special character handling
- Update scenarios

**Note**: Integration tests will skip gracefully if the OS keychain is not available (e.g., in headless CI environments).

### Running Tests

```bash
# Run all tests
go test ./internal/adapters/keychain/... -v

# Run with coverage
go test ./internal/adapters/keychain/... -cover

# Run benchmarks
go test ./internal/adapters/keychain/... -bench=.
```

### Test Results

```
=== Test Summary ===
Total Tests: 11
Passed: 4 (unit tests that don't require keychain)
Skipped: 7 (integration tests - keychain not available in environment)
Failed: 0

All tests pass successfully.
```

## Security Considerations

1. **No Credential Logging**: The implementation never logs usernames or API tokens
2. **OS-Level Encryption**: All credentials are encrypted by the OS keychain
3. **Validation**: All inputs are validated before storage
4. **Error Redaction**: Error messages do not expose sensitive data
5. **Idempotent Deletion**: Delete operations are safe to retry

## Error Handling

The adapter provides meaningful error messages for common scenarios:

- **Empty workspace ID**: "workspace ID cannot be empty"
- **Invalid configuration**: "invalid configuration: [specific field] is required"
- **Credential not found**: "credentials not found for workspace ID [id]"
- **Invalid service ID**: "invalid service ID: expected 'ticketr', got '[id]'"
- **Keychain errors**: "failed to [operation] credentials [in/from] keychain: [underlying error]"

## Cross-Platform Compatibility

### macOS
Uses the Keychain Access system (via Security framework).

### Windows
Uses the Credential Manager (via Windows Credential API).

### Linux
Uses the Secret Service API (requires `gnome-keyring` or `kwallet`).

**Important**: On Linux, the user's keychain must be unlocked for operations to succeed. In headless environments, tests will skip gracefully.

## Dependencies

- `github.com/zalando/go-keyring v0.2.6`
- Required OS dependencies:
  - macOS: Security framework (built-in)
  - Windows: Credential Manager (built-in)
  - Linux: `libsecret`, `gnome-keyring`, or `kwallet`

## File Structure

```
internal/adapters/keychain/
├── keychain_store.go           (269 lines) - Main implementation
├── keychain_store_test.go      (772 lines) - Integration tests
├── keychain_store_mock_test.go (229 lines) - Unit tests
├── interface_check.go          (6 lines)   - Compile-time interface check
└── README.md                               - This file
```

## Future Enhancements

Potential improvements for future iterations:

1. **Credential Rotation**: Automatic rotation of API tokens
2. **Backup/Restore**: Export/import encrypted credential backups
3. **Audit Logging**: Track credential access (without logging the actual credentials)
4. **Multi-tenancy**: Support for multiple users on the same machine
5. **Credential Expiry**: Time-based credential expiration

## License

This implementation is part of the Ticketr project and follows the project's license.
