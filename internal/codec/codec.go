package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	magicSize       = 4
	versionSize     = 4
	methodSize      = 32
	headerSize      = magicSize + versionSize + methodSize
	signatureSize   = sha256.Size
	resignatureSize = sha256.Size

	saltSize      = 16
	nonceSize     = 12
	keySize       = 32
	kdfIterations = 210_000

	version    uint32 = 1
	methodName        = "AES-256-GCM/PBKDF2-SHA256-V1"
)

var (
	magic           = [magicSize]byte{'Q', 'E', 'D', '!'}
	supportedMethod = fixedMethod(methodName)
	internalCode    = []byte("qed/internal-code/v1/CDU_Ge/2026")
)

func Encrypt(plaintext []byte, password string) ([]byte, error) {
	if password == "" {
		return nil, errors.New("password must not be empty")
	}

	header := makeHeader()

	salt, err := randomBytes(saltSize)
	if err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}
	nonce, err := randomBytes(nonceSize)
	if err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, header)
	body := make([]byte, 0, saltSize+nonceSize+len(ciphertext))
	body = append(body, salt...)
	body = append(body, nonce...)
	body = append(body, ciphertext...)

	packet := make([]byte, 0, headerSize+len(body)+signatureSize+resignatureSize)
	packet = append(packet, header...)
	packet = append(packet, body...)

	signature := sha256.Sum256(packet)
	packet = append(packet, signature[:]...)

	resignature := sha256.New()
	resignature.Write(packet)
	resignature.Write(internalCode)
	packet = resignature.Sum(packet)

	return packet, nil
}

func Decrypt(packet []byte, password string) ([]byte, error) {
	if password == "" {
		return nil, errors.New("password must not be empty")
	}

	body, header, err := parseAndVerify(packet)
	if err != nil {
		return nil, err
	}
	if len(body) < saltSize+nonceSize+16 {
		return nil, errors.New("qed body too short")
	}

	salt := body[:saltSize]
	nonce := body[saltSize : saltSize+nonceSize]
	ciphertext := body[saltSize+nonceSize:]

	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, header)
	if err != nil {
		return nil, errors.New("decrypt body: authentication failed")
	}
	return plaintext, nil
}

func makeHeader() []byte {
	header := make([]byte, headerSize)
	copy(header[:magicSize], magic[:])
	binary.BigEndian.PutUint32(header[magicSize:magicSize+versionSize], version)
	copy(header[magicSize+versionSize:headerSize], supportedMethod[:])
	return header
}

func parseAndVerify(packet []byte) ([]byte, []byte, error) {
	if len(packet) < headerSize+signatureSize+resignatureSize {
		return nil, nil, errors.New("qed packet too short")
	}
	if !bytes.Equal(packet[:magicSize], magic[:]) {
		return nil, nil, errors.New("invalid magic number")
	}
	if got := binary.BigEndian.Uint32(packet[magicSize : magicSize+versionSize]); got != version {
		return nil, nil, fmt.Errorf("unsupported version %d", got)
	}
	method := packet[magicSize+versionSize : headerSize]
	if !hmac.Equal(method, supportedMethod[:]) {
		return nil, nil, fmt.Errorf("unsupported method %q", trimMethod(method))
	}

	signatureStart := len(packet) - signatureSize - resignatureSize
	signatureEnd := signatureStart + signatureSize
	signature := packet[signatureStart:signatureEnd]
	resignature := packet[signatureEnd:]

	expectedSignature := sha256.Sum256(packet[:signatureStart])
	if !hmac.Equal(signature, expectedSignature[:]) {
		return nil, nil, errors.New("invalid signature")
	}

	hasher := sha256.New()
	hasher.Write(packet[:signatureEnd])
	hasher.Write(internalCode)
	if !hmac.Equal(resignature, hasher.Sum(nil)) {
		return nil, nil, errors.New("invalid re-signature")
	}

	header := packet[:headerSize]
	body := packet[headerSize:signatureStart]
	return body, header, nil
}

func deriveKey(password string, salt []byte) ([]byte, error) {
	key, err := pbkdf2.Key(sha256.New, password, salt, kdfIterations, keySize)
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}
	return key, nil
}

func randomBytes(size int) ([]byte, error) {
	buf := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func fixedMethod(name string) [methodSize]byte {
	var method [methodSize]byte
	copy(method[:], name)
	return method
}

func trimMethod(method []byte) string {
	if end := bytes.IndexByte(method, 0); end >= 0 {
		method = method[:end]
	}
	return string(method)
}
