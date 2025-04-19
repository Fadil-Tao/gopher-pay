package hashing

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type Argon2idHash struct {
	time uint32
	memory uint32 
	threads uint8
	keyLen uint32
	saltLen uint32
}

type HashSalt struct {
	Hash []byte
	Salt []byte
}


func NewArgon2idHash(time, saltLen uint32, memory uint32, threads uint8, keyLen uint32) *Argon2idHash {
	return &Argon2idHash{
		time: time,
		memory: memory,
		threads: threads,
		keyLen: keyLen,
		saltLen:  saltLen,
	}
} 

func randomSalt(length uint32) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (a *Argon2idHash) GenerateHash(password string, salt []byte)(*HashSalt, error){
	var err error 
	if len(salt) == 0{
		salt, err = randomSalt(a.saltLen)
	}
	if err != nil {
		return nil, err 
	}

	hash := argon2.IDKey([]byte(password), salt, a.time, a.memory, a.threads, a.keyLen)

	return &HashSalt{
		Hash: hash,
		Salt: salt,
	}, nil
}

func (a *Argon2idHash) Compare(hash , salt []byte, password string) (bool, error) {
	hashSalt, err := a.GenerateHash(password, salt)
	if err != nil {
		return false, err
	} 
	if !bytes.Equal(hash, hashSalt.Hash) {
		return false , nil
	}
	return true, nil
} 