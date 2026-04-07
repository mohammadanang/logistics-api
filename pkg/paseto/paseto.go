package paseto

import (
	"errors"
	"time"

	"aidanwoods.dev/go-paseto"
)

type TokenMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

func NewTokenMaker(secretKey string) (*TokenMaker, error) {
	if len(secretKey) != 32 {
		return nil, errors.New("secret key must be exactly 32 bytes")
	}

	key, err := paseto.V4SymmetricKeyFromBytes([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	return &TokenMaker{symmetricKey: key}, nil
}

// CreateToken membuat token baru untuk user/kurir
func (maker *TokenMaker) CreateToken(userID string, role string, duration time.Duration) (string, error) {
	token := paseto.NewToken()

	// Set Payload / Claims
	token.SetAudience("logistics-api")
	token.SetIssuer("system")
	token.SetExpiration(time.Now().Add(duration))
	token.SetString("user_id", userID)
	token.SetString("role", role)

	// Encrypt token (V4 Local)
	return token.V4Encrypt(maker.symmetricKey, nil), nil
}

// VerifyToken mengecek apakah token valid dan belum expired
func (maker *TokenMaker) VerifyToken(tokenString string) (*paseto.Token, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.ForAudience("logistics-api"))
	parser.AddRule(paseto.IssuedBy("system"))

	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, tokenString, nil)
	if err != nil {
		return nil, err // Akan error jika expired atau diubah isinya
	}

	return parsedToken, nil
}
