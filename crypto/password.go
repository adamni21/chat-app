package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordHasher interface {
	Generate(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

const Mega_Byte = 1024

type argonParams struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

type argon2Hasher struct {
	*argonParams
	saltLen uint32
}

func NewArgon2Hasher() *argon2Hasher {
	return &argon2Hasher{
		argonParams: &argonParams{
			time:    10,
			memory:  64 * Mega_Byte,
			threads: 2,
			keyLen:  64,
		},
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

// func (a *argon2Hasher) hash(hash, salt []byte, time, memory uint32, threads uint8) []byte {

// }

func buildArgonString(hash, salt []byte, time, memory uint32, threads uint8) string {
	b64Hash := base64.StdEncoding.EncodeToString(hash)
	b64Salt := base64.StdEncoding.EncodeToString(salt)

	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s", "argon2id", argon2.Version, memory, time, threads, b64Salt, b64Hash)
}

func parseArgonString(argonString string) (params *argonParams, hash, salt []byte, err error) {
	parts := strings.Split(argonString, "$")

	hash, err = base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decoding hash '%s': %w", hash, err)
	}
	salt, err = base64.StdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decoding salt '%s': %w", salt, err)
	}

	var memory uint32
	var time uint32
	var threads uint8
	count, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return nil, nil, nil, err
	}
	if count != 3 {
		fmt.Errorf("didn't parse all params from argonString: '%s'", argonString)
	}

	params = &argonParams{
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  uint32(len(hash)),
	}

	return
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
