package crypto

import (
	"bytes"
	"encoding/hex"
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

func (a *argon2Hasher) Generate(password string) (string, error) {
	salt, err := GenerateRandomBytes(a.saltLen)
	if err != nil {
		return "", fmt.Errorf("generating random salt: %w", err)
	}

	hash := a.hash([]byte(password), salt)

	return a.buildArgonString(hash, salt), nil
}

func (a *argon2Hasher) Verify(password, argonString string) (bool, error) {
	params, hash, salt, err := parseArgonString(argonString)
	if err != nil {
		return false, fmt.Errorf("parsing argonString '%s': %w", argonString, err)
	}

	match := bytes.Equal(hash, params.hash([]byte(password), salt))

	return match, nil
}

func (p *argonParams) hash(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, p.time, p.memory, p.threads, p.keyLen)
}

func (p *argonParams) buildArgonString(hash, salt []byte) string {
	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%x$%x", "argon2id", argon2.Version, p.memory, p.time, p.threads, salt, hash)
}

func parseArgonString(argonString string) (params *argonParams, hash, salt []byte, err error) {
	parts := strings.Split(argonString, "$")

	var memory uint32
	var time uint32
	var threads uint8
	count, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return nil, nil, nil, err
	}
	if count != 3 {
		return nil, nil, nil, fmt.Errorf("didn't parse all params from argonString: '%s'", argonString)
	}

	decodedHash, err := hex.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decoding hash '%s': %w", parts[5], err)
	}
	decodedSalt, err := hex.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decoding salt '%s': %w", parts[4], err)
	}

	hash = decodedHash
	salt = decodedSalt
	params = &argonParams{
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  uint32(len(hash)),
	}
	err = nil

	return
}
