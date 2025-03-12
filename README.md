# go-keystore

A secure and flexible keystore implementation for managing authentication tokens and cryptographic keys in Go applications.

## Features

- Secure storage of authentication tokens with expiration
- ECDSA private key management
- Configurable storage location
- Thread-safe operations
- Comprehensive error handling
- Full test coverage

## Installation

```bash
go get github.com/theblitlabs/go-keystore
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/theblitlabs/go-keystore"
)

func main() {
    // Create a new keystore with default configuration
    ks, err := keystore.NewKeystore(keystore.Config{})
    if err != nil {
        panic(err)
    }

    // Save an authentication token
    err = ks.SaveToken("your-auth-token")
    if err != nil {
        panic(err)
    }

    // Load the token
    token, err := ks.LoadToken()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Loaded token: %s\n", token)
}
```

### Custom Configuration

```go
ks, err := keystore.NewKeystore(keystore.Config{
    DirPath:  "/custom/path/to/keystore",
    FileName: "custom-keystore.json",
})
```

### Private Key Management

```go
// Save a private key
err = ks.SavePrivateKey("your-private-key-hex")
if err != nil {
    panic(err)
}

// Load the private key
privateKey, err := ks.LoadPrivateKey()
if err != nil {
    panic(err)
}
```

## Error Handling

The package provides specific error types for common scenarios:

- `ErrEmptyToken`: Returned when attempting to save an empty token
- `ErrNoKeystore`: Returned when no keystore file exists
- `ErrTokenExpired`: Returned when the stored token has expired
- `ErrInvalidToken`: Returned when the stored token is invalid
- `ErrNoPrivateKey`: Returned when no private key exists in the keystore

## Security

- Files are stored with 0600 permissions (user read/write only)
- Directories are created with 0700 permissions
- Tokens automatically expire after 1 hour (configurable)
- Private keys are validated before storage

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
