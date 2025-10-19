package keychain

import "github.com/karolswdev/ticktr/internal/core/ports"

// Compile-time check to ensure KeychainStore implements ports.CredentialStore interface.
var _ ports.CredentialStore = (*KeychainStore)(nil)
