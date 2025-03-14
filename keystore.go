package keystore

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// DefaultDirName is the default directory name for the keystore
	DefaultDirName = ".parity"

	// DefaultFileName is the default file name for the keystore
	DefaultFileName = "keystore.json"

	// DefaultFileMode is the default file permissions for the keystore file
	DefaultFileMode = 0600

	// DefaultDirMode is the default directory permissions for the keystore directory
	DefaultDirMode = 0700

	// TokenExpiryDuration is the duration after which a token is considered expired
	TokenExpiryDuration = 1 * time.Hour
)

var (
	ErrEmptyToken   = errors.New("token cannot be empty")
	ErrNoKeystore   = errors.New("no keystore found - please authenticate first")
	ErrTokenExpired = errors.New("token has expired - please re-authenticate")
	ErrInvalidToken = errors.New("invalid token found in keystore")
	ErrNoPrivateKey = errors.New("no private key found in keystore")
)

type Config struct {
	DirPath  string
	FileName string
}

type Store struct {
	AuthToken  string `json:"auth_token,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	config     Config
}

func NewKeystore(cfg Config) (*Store, error) {
	if cfg.DirPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		cfg.DirPath = filepath.Join(homeDir, DefaultDirName)
	}

	if cfg.FileName == "" {
		cfg.FileName = DefaultFileName
	}

	if err := os.MkdirAll(cfg.DirPath, DefaultDirMode); err != nil {
		return nil, fmt.Errorf("failed to create keystore directory: %w", err)
	}

	return &Store{config: cfg}, nil
}

func (s *Store) path() string {
	return filepath.Join(s.config.DirPath, s.config.FileName)
}

func (s *Store) SaveToken(token string) error {
	if token == "" {
		return ErrEmptyToken
	}

	s.AuthToken = token
	s.CreatedAt = time.Now().Unix()

	return s.save()
}

func (s *Store) LoadToken() (string, error) {
	if err := s.load(); err != nil {
		return "", err
	}

	if s.AuthToken == "" {
		return "", ErrInvalidToken
	}

	if time.Now().Unix()-s.CreatedAt > int64(TokenExpiryDuration.Seconds()) {
		return "", ErrTokenExpired
	}

	return s.AuthToken, nil
}

func (s *Store) SavePrivateKey(privateKeyHex string) error {
	if _, err := crypto.HexToECDSA(privateKeyHex); err != nil {
		return fmt.Errorf("invalid private key format: %w", err)
	}

	s.PrivateKey = privateKeyHex
	return s.save()
}

func (s *Store) LoadPrivateKey() (*ecdsa.PrivateKey, error) {
	if err := s.load(); err != nil {
		return nil, err
	}

	if s.PrivateKey == "" {
		return nil, ErrNoPrivateKey
	}

	return crypto.HexToECDSA(s.PrivateKey)
}

func (s *Store) GetPrivateKeyHex() (string, error) {
	if err := s.load(); err != nil {
		return "", err
	}

	if s.PrivateKey == "" {
		return "", ErrNoPrivateKey
	}

	return s.PrivateKey, nil
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal keystore: %w", err)
	}

	if err := os.WriteFile(s.path(), data, DefaultFileMode); err != nil {
		return fmt.Errorf("failed to write keystore file: %w", err)
	}

	return nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path())
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w at %s", ErrNoKeystore, s.path())
		}
		return fmt.Errorf("failed to read keystore: %w", err)
	}

	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("failed to parse keystore: %w", err)
	}

	return nil
}
