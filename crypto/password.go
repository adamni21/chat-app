package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type PasswordHasher interface {
	Generate(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

const Mega_Byte = 1024

type argon2Hasher struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func NewArgon2Hasher() *argon2Hasher {
	return &argon2Hasher{
		time:    10,
		memory:  64 * Mega_Byte,
		threads: 2,
		keyLen:  64,
		saltLen: 16,
	}
}

func (a *argon2Hasher) Generate(password []byte) (string, error) {
	salt, err := generateRandomBytes(a.saltLen)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)

	return buildArgonString(hash, salt, a.time, a.memory, a.threads), nil
}

func (a *argon2Hasher) Verify(password, hash []byte) (bool, error) {
	salt, err := generateRandomBytes(16)
	if err != nil {
		return false, err
	}

	argon2.IDKey(password, salt, a.time, a.memory, a.threads, a.keyLen)

	return true, nil
}

func buildArgonString(hash, salt []byte, time, memory uint32, threads uint8) string {
	b64Hash := base64.StdEncoding.EncodeToString(hash)
	b64Salt := base64.StdEncoding.EncodeToString(salt)

	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s", "argon2id", argon2.Version, memory, time, threads, b64Salt, b64Hash)
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
