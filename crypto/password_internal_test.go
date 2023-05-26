package crypto

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseArgonString(t *testing.T) {
	hash := []byte("abcde")
	salt := []byte("salt")

	params := argonParams{
		time:    10,
		memory:  64,
		threads: 2,
		keyLen:  uint32(len(hash)),
	}

	argonString := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", params.memory, params.time, params.threads, salt, hash)

	parsedParams, parsedHash, parsedSalt, err := parseArgonString(argonString)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(params, *parsedParams) {
		t.Errorf("got %+v want %+v", *parsedParams, params)
	}
	if !reflect.DeepEqual(hash, parsedHash) {
		t.Errorf("got %v want %v", parsedHash, hash)
	}
	if !reflect.DeepEqual(salt, parsedSalt) {
		t.Errorf("got %v want %v", parsedSalt, salt)
	}
}
